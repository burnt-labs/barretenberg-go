# Security Policy

This policy covers the Go bindings and proof verification wrappers for the
Barretenberg ZK proving library in this repository.

It supplements the
[organization-wide policy](https://github.com/burnt-labs/.github/blob/main/SECURITY.md),
which governs anything not addressed here.

## Reporting a Vulnerability

**Do not open a public GitHub issue for a security vulnerability.**

| Type of finding                  | How to report                                         |
| -------------------------------- | ----------------------------------------------------- |
| Security vulnerability           | Email [security@burnt.com](mailto:security@burnt.com)  |
| Non-sensitive or operational bug | Open a GitHub issue on this repository                 |

Include the type of vulnerability, affected version, steps to reproduce, impact,
how an attacker would exploit it, and any known mitigations.

We acknowledge receipt within **5 business days** and provide a triage decision
within **14 days**. Active exploitation, or confirmed attacker awareness of an
unpatched vulnerability, escalates the issue to Critical handling regardless of
its original classification.

## Scope

This repository provides Go bindings over the upstream Barretenberg library. The
security boundary that matters here is the **binding layer**: how proofs,
verification keys, and public inputs cross between Go and the underlying C
library, and whether a verification result can be influenced by attacker-supplied
data.

Findings of particular interest include:

- A proof or verification key that verifies when it should not
- Memory safety issues at the cgo boundary reachable from attacker-supplied input
- Incorrect handling of public inputs, serialization, or length prefixes that
  changes what a verification result actually attests to
- Verification returning success on malformed, truncated, or adversarially
  crafted input

## Proof of Concept Requirements

**Reports must include a proof of concept.** Severity is assessed on
demonstrated impact under real-world constraints, not theoretical worst-case
scenarios.

For verification-correctness findings, include the exact proof, verification
key, and public input bytes that produce the incorrect result, along with a
runnable Go test that demonstrates it. A description of a suspected weakness
without inputs that exhibit it is not actionable.

Where a finding is exploitable through a consumer of this library — for example,
on-chain verification — a proof of concept demonstrating that path carries
substantially more weight than one confined to this repository in isolation.

## Out of Scope

**Assets**

- **The upstream Barretenberg C library itself.** Vulnerabilities originating in
  Barretenberg are not eligible here and should be reported to that project.
  Only code originating in this repository is covered
- Other cryptographic dependencies and Go standard library components
- The chain node and its modules — see [`burnt-labs/xion`](https://github.com/burnt-labs/xion/blob/main/SECURITY.md)
- Smart contracts — see [`burnt-labs/contracts`](https://github.com/burnt-labs/contracts/blob/main/SECURITY.md)

**Vulnerability classes**

- Denial of service, including resource exhaustion from oversized or malformed
  proof input, where the impact is confined to the calling process
- Timing variation without a demonstrated key or witness recovery path
- Theoretical cryptographic weaknesses without a working proof of concept
  against this implementation
- Attacks where the attacker's cost to execute exceeds the demonstrable harm
- Findings already remediated in the current release
- Best practices and informational findings

## Severity Characterization

| Severity     | Description                                                                                                   |
| ------------ | --------------------------------------------------------------------------------------------------------------- |
| **CRITICAL** | Forged proof accepted as valid, where the finding is exploitable against a production consumer of this library   |
| **HIGH**     | Verification accepts a proof it should reject, or memory corruption at the cgo boundary reachable from attacker-supplied input |
| **MEDIUM**   | Incorrect verification behaviour requiring specific preconditions, or input handling that misrepresents what a verification result attests to |
| **LOW**      | Valid, reproducible code-level issue with no direct impact on verification soundness, representing a meaningful hardening opportunity |

Severity is assessed by Burnt Labs based on demonstrated impact. Reports
submitted at a severity that does not match the definitions above are assessed
as written; we do not reclassify or negotiate severity on a reporter's behalf.

## Responsible Disclosure

- Do not exploit a vulnerability beyond what is necessary to confirm it exists
- **Do not test against XION mainnet.** Testing that targets live production
  systems will disqualify the report
- Do not access, modify, or exfiltrate user data
- Do not disclose publicly before a fix is confirmed and deployed

## Safe Harbor

Burnt Labs will not pursue legal action against researchers who report
vulnerabilities in good faith under this policy, do not exploit beyond what is
necessary to confirm the finding, do not access or disclose user data, and do
not disrupt production systems.

Authorization to actively test extends only to assets named in a published Burnt
Labs bug bounty program. Testing systems outside that scope is not authorized.
Reporting a vulnerability you encountered incidentally is always welcome.
