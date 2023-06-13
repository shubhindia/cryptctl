# `cryptctl` CLI Tool for managing [EncryptedSecrets](https://github.com/shubhindia/encrypted-secrets)

## What is `cryptctl`?

**cryptctl** is a simple command-line interface (CLI) tool designed to facilitate the management of EncryptedSecrets.
With Cryptctl, you can easily update encrypted secrets within your Kubernetes cluster, ensuring the secure handling of sensitive information.

### Features
**Effortless Encryption:** Cryptctl simplifies the process of encrypting secrets by providing a straightforward command-line interface. It handles the encryption and decryption operations seamlessly, making it easy to work with encrypted secrets in your Kubernetes environments.

**Simplified Management:** Since, the secrets are encrypted, they can be easily stored in a repository. Once, the `EncryptedSecret` object is applied, `encrypted-secrets` controller takes care of decrypting the provided secrets and creates a k8s secret with decrpted values. Which can be access by the pod as required.

https://github.com/shubhindia/cryptctl/assets/7694806/c36853b8-e373-4c38-840c-0ff9a8b9cfcf


