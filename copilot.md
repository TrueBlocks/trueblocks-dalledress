# Copilot Configuration

- Never use `npm`. Always use `yarn` for package management.

- Don't build from the `./frontend` folder. Always just do `yarn build` from the root of the repo.

- Never try to run the app from the `./frontend` folder. If you must run the app, do so from the root of the repo with `yarn start`.

- Generally, don't run tests, but if you must, run `yarn test` from the root of the repo.

- Never use console.log in the frontend. Instead, use Log from the @utils. This takes a single string, but it prints out to the console. Must easier to use.

- When generating Wails bindings, use `wails generate module` from the root of the repo.