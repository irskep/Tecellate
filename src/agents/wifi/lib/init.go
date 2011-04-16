package lib

import pseudo_rand "rand"
import crypto_rand "crypto/rand"
import . "byteslice"

// initializer for random number generator -------------------------------------
func randbytes(k int) ByteSlice {
    bytes := make(ByteSlice, k)
    cbytes := bytes[:]
    for
        n, err := crypto_rand.Read(cbytes);
        n < k;
        n, err = crypto_rand.Read(cbytes) {
            if err != nil {
                panic("Can't get random bytes.")
            }
            k = k-n
            cbytes = cbytes[n:]
    }
    return bytes
}
func init() {
    pseudo_rand.Seed(int64(randbytes(8).Int64()))
}
