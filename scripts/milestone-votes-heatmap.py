#!/usr/bin/env python3
"""Render a validator x milestone participation heatmap from the JSON
output of `polycli heimdall milestone votes --json`.

Each column is one finalized milestone; each row is one validator. The
cell shows what that validator's vote looked like at the height the
milestone finalized from:

    covered   the validator's proposition covered the milestone end block
    lag 1     present but one bor block behind the 2/3 majority
    lag 2+    present but two or more bor blocks behind
    no prop   signed the commit but carried no milestone proposition
    absent    never delivered a pre-commit at that height

Horizontal streaks of color identify a single broken validator;
vertical bands identify network-wide events (usually bor, not
heimdall). The healthy state is deliberately near-white so anomalies
are the only ink that draws attention.

Row labels include the validator's registered name when available.
Names are not stored on heimdall or in the L1 staking contracts; they
are profile data served by Polygon's staking API, so the script
fetches them from there (best-effort: a missing or unreachable API
just leaves names blank). --names picks the network, with auto
matching the data's signers against both networks.

Usage:
    polycli heimdall milestone votes --json > votes.json
    scripts/milestone-votes-heatmap.py votes.json
    scripts/milestone-votes-heatmap.py votes.json -o heatmap.png --sort misses
    scripts/milestone-votes-heatmap.py votes.json --names off
    polycli heimdall milestone votes --json | scripts/milestone-votes-heatmap.py -

Requires matplotlib and numpy (pip install matplotlib numpy).
"""

import argparse
import json
import sys
import urllib.error
import urllib.request

import numpy as np
import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt
from matplotlib.colors import ListedColormap
from matplotlib.patches import Patch

# Cell states, ordered from healthy to worst. Index doubles as the
# matrix value; NO_DATA is masked out (rendered white).
COVERED, LAG_1, LAG_2_PLUS, NO_PROP, ABSENT, NO_DATA = range(6)

STATE_LABELS = {
    COVERED: "covered",
    LAG_1: "lag 1",
    LAG_2_PLUS: "lag 2+",
    NO_PROP: "commit, no proposition",
    ABSENT: "absent",
}

STATE_COLORS = {
    COVERED: "#e9ece6",  # near-white: healthy cells should be quiet
    LAG_1: "#f5c842",
    LAG_2_PLUS: "#e8743b",
    NO_PROP: "#8e6bb8",
    ABSENT: "#c0392b",
}


# Validator display names are off-chain profile data served by
# Polygon's staking API (they exist neither on heimdall nor in the L1
# staking contracts). Keyed by lowercase signer address.
STAKING_API_URLS = {
    "mainnet": "https://staking-api.polygon.technology/api/v2/validators?limit=1000",
    "amoy": "https://staking-api-amoy.polygon.technology/api/v2/validators?limit=1000",
}


def fetch_names(network):
    """Return {lowercase signer: name} from the staking API, or {} on
    any failure — names are decoration, never worth failing the run."""
    url = STAKING_API_URLS[network]
    # The staking API rejects urllib's default User-Agent with a 403.
    req = urllib.request.Request(url, headers={"User-Agent": "milestone-votes-heatmap/1.0"})
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            payload = json.load(resp)
    except (urllib.error.URLError, TimeoutError, json.JSONDecodeError, OSError) as e:
        print(f"warning: could not fetch {network} validator names: {e}", file=sys.stderr)
        return {}
    names = {}
    for v in payload.get("result", []):
        signer = str(v.get("signer", "")).lower()
        name = str(v.get("name", "")).strip()
        if signer and name:
            names[signer] = name
    return names


def resolve_names(mode, data):
    """Build the signer->name map for --names MODE.

    auto fetches both networks and keeps whichever matches more of the
    signers present in the data — this also catches captures made
    against a node whose REST/RPC endpoints were cross-wired.
    """
    if mode == "off":
        return {}
    if mode != "auto":
        return fetch_names(mode)
    signers = {v["signer"].lower() for v in data["votes"]}
    best, best_hits = {}, 0
    for network in STAKING_API_URLS:
        names = fetch_names(network)
        hits = sum(1 for s in signers if s in names)
        if hits > best_hits:
            best, best_hits = names, hits
    if not best_hits:
        print("warning: no validator names matched this data on any network", file=sys.stderr)
    return best


def classify(vote):
    """Map one vote record to a cell state."""
    if vote.get("flag") != "COMMIT":
        return ABSENT
    if vote.get("prop_start") is None:
        return NO_PROP
    milestone = vote.get("milestone")
    if milestone == "miss":
        lag = vote.get("lag")
        # A miss without a computable lag still means "behind".
        if lag is None or lag >= 2:
            return LAG_2_PLUS
        return LAG_1
    # Covered (milestone == its number). A null milestone cannot occur
    # here because columns are built only from finalize heights, but
    # treat it as covered-equivalent rather than inventing a state.
    return COVERED


def short_addr(signer):
    return signer[:6] + "…" + signer[-4:] if len(signer) > 12 else signer


def build_matrix(data):
    """Pivot votes into (matrix, row_meta, milestones).

    Rows are validators, columns are milestones ordered by number.
    """
    milestones = sorted(data.get("milestones", []), key=lambda m: int(m["number"]))
    if not milestones:
        sys.exit("input contains no finalized milestones; nothing to plot")

    # vote_height -> list of column indexes (rarely more than one).
    height_cols = {}
    for col, ms in enumerate(milestones):
        height_cols.setdefault(ms["vote_height"], []).append(col)

    # signer -> row index, plus presentation metadata.
    rows = {}
    row_meta = []
    for vote in data["votes"]:
        signer = vote["signer"]
        if signer not in rows:
            rows[signer] = len(row_meta)
            row_meta.append(
                {
                    "signer": signer,
                    "val_id": vote.get("val_id", "-"),
                    "power": vote.get("power", 0),
                    "problems": 0,
                }
            )

    matrix = np.full((len(row_meta), len(milestones)), NO_DATA, dtype=np.int8)
    for vote in data["votes"]:
        cols = height_cols.get(vote["height"])
        if not cols:
            continue  # no milestone finalized from this height
        state = classify(vote)
        row = rows[vote["signer"]]
        for col in cols:
            matrix[row, col] = state
        if state != COVERED:
            row_meta[rows[vote["signer"]]]["problems"] += len(cols)
    return matrix, row_meta, milestones


def sort_rows(matrix, row_meta, mode):
    if mode == "misses":
        # Worst rows render first (top of the image); ties broken by
        # power so ordering is stable.
        key = lambda i: (row_meta[i]["problems"], row_meta[i]["power"])
    else:
        key = lambda i: row_meta[i]["power"]
    order = sorted(range(len(row_meta)), key=key, reverse=True)
    return matrix[order], [row_meta[i] for i in order]


def time_span(data, milestones):
    """First/last vote time of the plotted milestone heights."""
    heights = {ms["vote_height"] for ms in milestones}
    times = sorted(v["time"] for v in data["votes"] if v["height"] in heights)
    if not times:
        return ""
    trim = lambda t: t.split(".")[0].replace("T", " ").rstrip("Z") + "Z"
    return f"{trim(times[0])} → {trim(times[-1])}"


def render(matrix, row_meta, milestones, data, output, dpi):
    n_rows, n_cols = matrix.shape
    fig_w = max(10, min(24, n_cols * 0.035))
    fig_h = max(4, n_rows * 0.32 + 1.8)
    fig, ax = plt.subplots(figsize=(fig_w, fig_h))

    cmap = ListedColormap([STATE_COLORS[s] for s in range(ABSENT + 1)])
    masked = np.ma.masked_equal(matrix, NO_DATA)
    ax.imshow(
        masked,
        aspect="auto",
        interpolation="nearest",
        cmap=cmap,
        vmin=COVERED,
        vmax=ABSENT,
    )

    name_w = min(20, max((len(m["name"]) for m in row_meta), default=0))
    labels = []
    for m in row_meta:
        name = m["name"][:name_w]
        label = f"val {m['val_id']:>3}  "
        if name_w:
            label += f"{name:<{name_w}}  "
        labels.append(label + short_addr(m["signer"]))
    ax.set_yticks(range(n_rows))
    ax.set_yticklabels(labels, fontsize=8, fontfamily="monospace")

    tick_step = max(1, n_cols // 12)
    ticks = list(range(0, n_cols, tick_step))
    ax.set_xticks(ticks)
    ax.set_xticklabels(
        [milestones[i]["number"] for i in ticks], fontsize=8, rotation=45, ha="right"
    )
    ax.set_xlabel("milestone number")

    span = time_span(data, milestones)
    ax.set_title(
        f"Milestone vote participation — {n_rows} validators × "
        f"{n_cols} milestones" + (f"\n{span}" if span else ""),
        fontsize=11,
    )

    counts = {s: int((matrix == s).sum()) for s in STATE_LABELS}
    ax.legend(
        handles=[
            Patch(facecolor=STATE_COLORS[s], label=f"{STATE_LABELS[s]} ({counts[s]})")
            for s in STATE_LABELS
        ],
        loc="upper center",
        bbox_to_anchor=(0.5, -0.18),
        ncol=len(STATE_LABELS),
        fontsize=8,
        frameon=False,
    )

    fig.tight_layout()
    fig.savefig(output, dpi=dpi, bbox_inches="tight")
    print(f"wrote {output} ({n_rows} validators × {n_cols} milestones)")


def print_problem_summary(row_meta, n_cols):
    problems = [m for m in row_meta if m["problems"] > 0]
    if not problems:
        print("no problem cells: every validator covered every milestone")
        return
    problems.sort(key=lambda m: m["problems"], reverse=True)
    name_w = max((len(m["name"]) for m in problems[:10]), default=0)
    print("validators with problem cells (not covered / behind / absent):")
    for m in problems[:10]:
        pct = 100.0 * m["problems"] / n_cols
        name = f"{m['name']:<{name_w}}  " if name_w else ""
        print(
            f"  val {m['val_id']:>4}  {name}{m['signer']}  "
            f"{m['problems']:>5}/{n_cols} ({pct:.1f}%)"
        )


def main():
    parser = argparse.ArgumentParser(
        description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter
    )
    parser.add_argument(
        "input", help="JSON file from `polycli heimdall milestone votes --json`, or - for stdin"
    )
    parser.add_argument(
        "-o", "--output", help="output image path (default: <input>.png, or milestone-votes.png for stdin)"
    )
    parser.add_argument(
        "--sort",
        choices=["power", "misses"],
        default="power",
        help="row order: voting power (stable identity) or problem count (worst rows on top)",
    )
    parser.add_argument(
        "--names",
        choices=["auto", "amoy", "mainnet", "off"],
        default="auto",
        help="fetch validator display names from the Polygon staking API "
        "(auto picks the network whose validators match the data)",
    )
    parser.add_argument("--dpi", type=int, default=150, help="output image dpi")
    args = parser.parse_args()

    if args.input == "-":
        data = json.load(sys.stdin)
        output = args.output or "milestone-votes.png"
    else:
        with open(args.input) as f:
            data = json.load(f)
        output = args.output or args.input.rsplit(".json", 1)[0] + ".png"

    if "votes" not in data:
        sys.exit("input does not look like `milestone votes --json` output (no .votes key)")

    names = resolve_names(args.names, data)
    matrix, row_meta, milestones = build_matrix(data)
    for m in row_meta:
        name = names.get(m["signer"].lower(), "")
        if not name and names:
            # Staking-dashboard convention for unregistered validators.
            name = "Anonymous" if m["val_id"] == "-" else f"Anonymous {m['val_id']}"
        m["name"] = name
    matrix, row_meta = sort_rows(matrix, row_meta, args.sort)
    render(matrix, row_meta, milestones, data, output, args.dpi)
    print_problem_summary(row_meta, matrix.shape[1])


if __name__ == "__main__":
    main()
