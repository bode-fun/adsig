# TODO

## Configuration

- [ ] Get config from standard unix config locations.
- [ ] Parse environment variables from YAML config.
- [ ] Validate the config values.

## Templates

- [ ] Store generated templates in the filesystem under `/var/lib/adsig`, because `/var/tmp` might be deleted. [Ref](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch05s08.html)
- [ ] Add a fs-lock to prevent that a file is written to and read from at the same time.

## API

- [ ] Add API to request signatures for a given user.

## Optimizations

- [ ] Don't rebuild the Group/Member tree every request.
- [ ] Use a cache or DB for the generated Signatures.
