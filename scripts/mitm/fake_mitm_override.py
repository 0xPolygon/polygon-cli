# fake_mitm_override_fixed.py
# this is a mitm proxy script to be used to fake certain JSON-RPC responses from an Ethereum node
import json, time
from mitmproxy import http

def _serialize_params(p):
    if p is None:
        return None
    return json.dumps(p, separators=(",", ":"), sort_keys=True)

# keys: (method, params_json_or_None)
OVERRIDES = {
    ("txpool_content", None): {"error": {"code": -32601, "message": "Method not found"}},
    # ("txpool_content", None): {"jsonrpc": "2.0","id": 1,"result": {"pending": {"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266": {"40": {}}}},
    # ("eth_getTransactionCount", _serialize_params(["0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "latest"])): {"jsonrpc": "2.0", "id": 1, "result": "0xa"},
    ("eth_getTransactionCount", _serialize_params(["0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "pending"])): {"jsonrpc": "2.0", "id": 1, "result": "0x100"},
}

def request(flow: http.HTTPFlow):
    body_text = flow.request.get_text()
    try:
        j = json.loads(body_text)
        method = j.get("method")
        params = j.get("params")
        req_id = j.get("id")
        print(f"[{time.time():.3f}] REQ method={method} id={req_id} params={_serialize_params(params)}")
    except Exception:
        # malformed JSON: log truncated raw body and let it pass
        print(f"[{time.time():.3f}] REQ malformed body: {body_text[:200]!r}")
        return

    key_exact = (method, _serialize_params(params))
    key_any = (method, None)
    override = OVERRIDES.get(key_exact) or OVERRIDES.get(key_any)
    if override is not None:
        resp = {"jsonrpc": "2.0", "id": req_id}
        resp.update(override)
        flow.response = http.Response.make(200, json.dumps(resp), {"Content-Type": "application/json"})
