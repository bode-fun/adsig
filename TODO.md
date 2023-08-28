# TODO

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

- [ ] Don't rebuild the Group/Member tree every request.
- [ ] Try to reduce the references kept around in the Group/Member tree.
