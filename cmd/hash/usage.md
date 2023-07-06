The `hash` command provides a simple mechanism to perform hashes on files, standard input, and arguments. Below shows various ways to provide input.

```bash
$ echo -n "hello" > hello.txt
$ polycli hash sha1 --file hello.txt
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
$ echo -n "hello" | polycli hash sha1
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
$ polycli hash sha1 hello
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
```

We've provided many standard hashing functions.

```bash
echo -n "hello" | polycli hash md4
echo -n "hello" | polycli hash md5
echo -n "hello" | polycli hash sha1
echo -n "hello" | polycli hash sha224
echo -n "hello" | polycli hash sha256
echo -n "hello" | polycli hash sha384
echo -n "hello" | polycli hash sha512
echo -n "hello" | polycli hash ripemd160
echo -n "hello" | polycli hash sha3_224
echo -n "hello" | polycli hash sha3_256
echo -n "hello" | polycli hash sha3_384
echo -n "hello" | polycli hash sha3_512
echo -n "hello" | polycli hash sha512_224
echo -n "hello" | polycli hash sha512_256
echo -n "hello" | polycli hash blake2s_256
echo -n "hello" | polycli hash blake2b_256
echo -n "hello" | polycli hash blake2b_384
echo -n "hello" | polycli hash blake2b_512
echo -n "hello" | polycli hash keccak256
echo -n "hello" | polycli hash keccak512
```
