# Poor's man debugger

<h3>🚧 Under development, it works but some stuff might change in the future.</h3>

<hr>

The most simple way to debug your code.

This is a simple TUI debugger, it starts a local server that accepts payloads from the adapters. Everything goes through HTTP requests, so it's very fast compared to other solutions.
You can think of it as a console.log(), but instead of the browser you use the terminal and it works with any language.

<img width="908" alt="_zsh_tmux_plugin_run 2022-04-29 15-28-12" src="https://user-images.githubusercontent.com/35064680/165953816-470e13b6-fde5-412e-b7e7-89c331a2d284.png">


## Requirements

- PMD-adapter for your language installed in the project.

## Installation

TODO

## Configuration

TODO

## Usage

TODO

### Keybindings

TODO

### Adapter API

An example API for an adapter

```
curl --request POST \
  --url http://localhost:3000/dump \
  --header 'Content-Type: application/json' \
  --data '{
	"timestamp": "3223232",
	"line": "6",
	"connector_type": "php",
	"filepath": "/home/nkoporec/personal/drupal/web/index.php",
	"callstack": [
		"10:/home/nkoporec/personal/drupal/web/index.php",
		"100:/home/nkoporec/personal/drupal/web/index.php",
		"2213:/home/nkoporec/personal/drupal/web/index.php"
	],
	"payload": ""
}'
```

Types:
 - timestamp -> String
 - line -> String
 - connector_type -> String
 - filepath -> String
 - callstack -> Array (where key is a line number (int), and value is a file path (string)
 - payload -> JSON encoded string

## Changelog

Please see [CHANGELOG](CHANGELOG.md) for more information on what has changed recently.

## Contributing

Please see [CONTRIBUTING](.github/CONTRIBUTING.md) for details.

## Security Vulnerabilities

Please review [our security policy](../../security/policy) on how to report security vulnerabilities.

## Credits

- [nkoporec](https://github.com/nkoporec)
- [All Contributors](../../contributors)

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
