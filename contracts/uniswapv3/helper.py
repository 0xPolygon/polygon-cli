import math
import sys

def price_to_tick(p):
  return math.floor(math.log(p, 1.0001))

def price_to_sqrtp(p):
  q96 = 2**96
  return int(math.sqrt(p) * q96)

def main():
  if len(sys.argv) != 3:
    print("Usage: python3 helper.py <reserve_x> <reserve_y>")
    sys.exit(1)

  try:
    reserve_x = int(sys.argv[1])
    reserve_y = int(sys.argv[2])
  except ValueError:
    print("Invalid input. Please provide numeric values for reserve_x and reserve_y.")
    sys.exit(1)

  price = int(reserve_y / reserve_x)
  print("Current price:", price)
  print("Current price (Q64.96):", price_to_sqrtp(price))
  print("Tick index:", price_to_tick(price))

if __name__ == "__main__":
  main()
