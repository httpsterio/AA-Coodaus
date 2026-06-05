# AA-coodaus VS Code extension

Syntax highlighting for `.aa` files.

## Features

- Keywords, strings, numbers, booleans, operators
- String interpolation `{muuttuja}`
- Comments (`--`) and shebang lines
- Escape sequences (`\n`, `\t`, `\\`, `\"`)

## Install

```sh
cp -r plugins/vscode ~/.vscode/extensions/aakoodaus-0.1.0/
```

Then reload VS Code (`Ctrl+Shift+P` → "Developer: Reload Window").

## Package

```sh
cd plugins/vscode
npx vsce package
```

Produces `aakoodaus-0.1.0.vsix`.
