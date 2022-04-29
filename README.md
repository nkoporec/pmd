# Poor's man debugger

<h3>🚧 Under development, it works but some stuff might change in the future.</h3>

<hr>

The most simple way to debug your code.

This is a simple TUI debugger, it starts a local server that accepts payloads from the adapters. Everything goes through HTTP requests, so it's very fast compared to other solutions.
You can think of it as a console.log(), but instead of the browser you use the terminal and it works with any language.

<img width="1018" alt="_zsh_tmux_plugin_run 2022-04-29 13-27-39" src="https://user-images.githubusercontent.com/35064680/165936182-3717ddda-d380-40cf-b17d-37b9f85e41ad.png">


## Requirements

- PMD-adapter for your language installed in the project.

## Installation

You can install the package via go install:

```bash
go install github.com/nkoporec/pmd@latest
```

Install one of the supported adapters:

 - [PMD-PHP](https://www.github.com/nkoporec/pmd-php)

## Configuration

By default it will start a server at 127.0.0.1:8080. If this port does not work for you, then you can change it in the config file, located at `$HOME/.config/pmd/config.yml`

## Usage

PHP example:

1. Start the PMD with `pmd listen`
2. Set the breakpoint

```php
pmd("Hello world!");
```

3. Run the code and check the terminal


### Commands

- Move up -> `K`
- Move down -> `J`
- Clear screen -> `Ctrl + r`

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
