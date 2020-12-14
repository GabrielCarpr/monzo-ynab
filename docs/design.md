# Monzo -> YNAB integration

## Problem
Currently, YNAB only has native integrations with US banks. For UK residents, this is annoying.

To bypass this problem natively, you have to manually add transactions, which is slow, annoying, and causes your
budget to only be up to date immediately after performing reconciliation. This isn't enough to feel confident about
your spending and budget.

## Goals
Make an easy integration to allow Monzo account holders to easily import their transactions in real time.

Minorly technical users should be able to easily run the tool as a one off, or set it up to process transactions in real time,
without having to deal with a complicated deployment.

## Current solutions

### Manual imports
- Slow
- Forget to do them
- Often create inaccuracies by missing a penny

### YNAB for fintech
- Either a node or Ruby app, which is difficult to deploy, even with instructions. Requires many dependencies, and updating is even more difficult

### Sync for YNAB
- Expensive
- No visibility into sync status and transactions

## Solution
A CLI/API app, distributed as a single Go binary. CLI installs, configures, and deploys without any intervention.

CLI can be used for a one time sync, or to start a webhook server which accepts webhooks from Monzo.

When running as an app, or even locally, a web interface can be started which will offer more insight into transactions (premium only?).

### Example installation flow

```
./monzo-ynab install
What's your YNAB personal API token?
> 

```

### CLI options
```
install
run [--port]
set-option <option> <value>
sync <previous days>
update [--enable-auto]
```

## Implementation
- Use PromptUI for CLI, which easily creates an interactive console application


## Questions
- Will YNAB dedupe transactions automatically?
- How do I sync a bunch of transactions?
- On receipt of a webhook from Monzo, how do I verify it's validity?
- Will I need to request other information from Monzo before sending it to YNAB?
- Will I need to perform any other operations on YNAB before sending the transaction?
    - As long as I have the budget ID, and the account ID, I can use all the information from the transaction in Monzo to create the YNAB transaction.
- How do I let users authenticate with Monzo?