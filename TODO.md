# TODO

## General

- [ ] Modify the registry to set the default signature.
- [ ] Get the sAMAccountname from Windows directly.

## Configuration

- [ ] Get config from standard unix config locations or remote server.
- [ ] Parse environment variables from YAML config.
- [ ] Validate the config values.

## Templates

- [ ] Store and render templates in memory
  - [ ] Save signatures to `%APPDATA%\Microsoft\Signatures`

## Error Handling

- [ ] Add multi-error, where it provides usable information.
- [ ] Wrap errors.
- [ ] Recover from errors, where it makes sense.

## Optimizations

- [ ] Try to reduce the references kept around in the Group/Member tree.

## Daemon

- [ ] Add a daemon mode for small companies, which usually don't use a device management solution.

## Testing

- [ ] Write unit tests
- [ ] Setup a test AD
