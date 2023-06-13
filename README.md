# `cryptctl` CLI Tool for managing [EncryptedSecrets](https://github.com/shubhindia/encrypted-secrets)

## Note: `cryptctl` is currently a work in progress and is in the alpha stage. Please use it with caution in production environments.
## What is `cryptctl`?

**cryptctl** is a simple command-line interface (CLI) tool designed to facilitate the management of EncryptedSecrets.
With Cryptctl, you can easily update encrypted secrets within your Kubernetes cluster, ensuring the secure handling of sensitive information.

### Features
- **Effortless Encryption:** Cryptctl simplifies the process of encrypting secrets by providing a straightforward command-line interface. It handles the encryption and decryption operations seamlessly, making it easy to work with encrypted secrets in your Kubernetes environments.

- **Simplified Management:** Since, the secrets are encrypted, they can be easily stored in a repository. Once, the `EncryptedSecret` object is applied, `encrypted-secrets` controller takes care of decrypting the provided secrets and creates a k8s secret with decrpted values. Which can be access by the pod as required.


Here's a **`cryptctl`** demo:

![cryptctl demo GIF](img/cryptctl-edit-demo.gif)
