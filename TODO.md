# TODO

## Configuration

- [ ] Get config from standard unix config locations.
- [ ] Parse environment variables from YAML config.
- [ ] Validate the config values.

## Templates

- [ ] Store DB in the filesystem under `/var/lib/adsig`, because `/var/tmp` might be deleted. [Ref](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch05s08.html)

## API

- [ ] Add API to request signatures for a given user.
  - [ ] Use DN or sAMAccountName as identifier, because multiple user might share an email address.

## Error Handling

- [ ] Add multi-error, where it provides usable information.
- [ ] Wrap errors.
- [ ] Recover from errors, where it makes sense.

## Optimizations

- [ ] Don't rebuild the Group/Member tree every request.
- [ ] Try to reduce the references kept around in the Group/Member tree.
