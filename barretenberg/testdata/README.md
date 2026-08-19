# Test Vectors

Binary test vectors for UltraHonk proof verification, generated with Aztec Barretenberg v5.2.0.

## Files

- `statics/vk` — Verification key (binary, UltraHonk format)
- `statics/proof` — Proof (binary)
- `statics/public_inputs` — Concatenated 32-byte field elements (big-endian)

Vectors are flavor-specific: a proof produced by one Barretenberg release does
not verify against another release whose transcript or flavor changed. Whenever
`checksums.json` moves to a new `aztec_tag`, regenerate all three files together
and confirm `TestVerifyValidProof` passes before releasing.

## Regenerating

Requires [Noir](https://noir-lang.org/) (`nargo` 1.0.0-beta.26) and the
[Barretenberg](https://github.com/AztecProtocol/aztec-packages) CLI (`bb`)
matching the `aztec_tag` in `checksums.json`. The `bb` CLI ships as
`barretenberg-<arch>-<os>.tar.gz` on the aztec-packages release.

The vectors come from the stock Noir starter circuit with one public input:

```noir
// src/main.nr
fn main(x: Field, y: pub Field) {
    assert(x != y);
}
```

```toml
# Prover.toml
x = "1"
y = "5"
```

```bash
nargo execute witness
bb prove \
    -b target/<package>.json \
    -w target/witness.gz \
    -o out \
    --write_vk \
    --verify
cp out/proof out/vk out/public_inputs testdata/statics/
```

`bb prove` defaults to the `ultra_honk` scheme with ZK enabled, which is the
flavor `wrapper/barretenberg_wrapper.cpp` verifies (`bb::UltraZKFlavor`). Do not
pass `--verifier_target`; the non-default targets change the transcript hash and
produce proofs this wrapper rejects.
