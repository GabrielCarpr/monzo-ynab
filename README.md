# Monzo-to-YNAB

This lightweight app allows you to connect Monzo to YNAB to automatically sync transactions.

You can run monzo-ynab as a webhook server, or bulk sync Monzo to YNAB from your computer.

## Installation

1. Put `monzo-ynab` in your path
2. Run `monzo-ynab install`. An interactive prompt will take your credentials and save it to a file (`/etc/monzo-ynab/config.json`). Alternatively it can read from environment variables, however you'll have to read the source to find out the names (`src/internal/config`).
3. You can now setup Monzo webhooks by running `monzo-ynab register-webhooks`.
4. You can sync the past x days of transactions by running `monzo-ynab sync [days]`. This is idempotent, so feel free to run it multiple times if you want to test your config.
5. You can run the webhook server by running `monzo-ynab run [port]`

## Other things
- If you get your Monzo access token from Monzo developer portal (which you most likely have), it will only last for ~24 hours. Webhooks will continue to function after that, but sync will fail. If you want to update your Monzo Access token, you can run `monzo-ynab configure MonzoAccessToken [token]`.
