{
  function testADD(loops) -> result
  {
    result := 1
    for { let i := 0 } lt(i, loops) { i := add(i, 1) }
    {
      result := add(result, 1)
    }
  }
}
