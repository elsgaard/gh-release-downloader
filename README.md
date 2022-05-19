# GitHub release downloader

Go based GitHub release downloader CLI

## Installation

- Fedora/CentOS/RedHat

```bash
export GHRDVER=1.1.2 && sudo curl -L https://github.com/elsgaard/gh-release-downloader/releases/download/v$GHRDVER/ghrd -o /usr/local/bin/ghrd \
&& sudo chmod +x /usr/local/bin/ghrd \
&& sudo ln -s /usr/local/bin/ghrd /usr/bin/ghrd \
&& unset GHRDVER
```

## Usage

```bash
Usage ./ghrd
  -a string
        Artifact name (no default)
  -o string
        path/filename (no default)
  -p string
        Personal access token (no default)
  -r string
        A release version (default: 'latest') (default "latest")
  -repo string
        Repository like elsgaard/gh-release-download
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

This project is licensed under the Apache License - see the [LICENSE](./LICENSE) file for details

## References

Inspired by the bash version made by https://github.com/zero88/gh-release-downloader
