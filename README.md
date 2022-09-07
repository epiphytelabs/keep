# keep

## Installation

```
brew tap epiphytelabs/tap
brew install keep
keep server install
```

## Install an app

```
keep install firefly
```

## (Optional) Trust the `*.app.keep` certificate

```
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/.keep/ssl/cert.pem
```
