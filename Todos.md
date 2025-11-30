# Essential Features

## Scrollable Content - Easy

Content needs to be scrollable both by shift down, shift up, J, K and pgDwn, pgUp, and D and U

## Copiable Content - Easy

Any "Main Content" should have easy to copy-to-clipboard contents with "c"

## Raw JSON - Medium

Any resource should be viewable / expandable / editable by JSON or YAML funcitonality

## Events Pane - Easy

Events need to have a hotkey (likely e) for selected apps to see useful event information

## Pod Logs Pane - Medium

Pod logs need to be viewable in a pane with a shortcut (likely l)

## Applications Pane Actions - Easy

ArgoCD applications should have at least these actions:

1. Sync: hotkey (s) (async with loading animation)
2. Refresh: hotkey (r) (async with loading animation)
3. Edit: hotkey(e)? manifest in editor (maybe refine edits for target revisions etc. as time goes on)

## Global Default / User Config - Unknown

There should be a global config filled with different defaults that can be overriden by a user.

- Use .config/argocd-tui/config.yaml structure
- Configurable theme
- pagination limits for different panes
- default editor

## Help menu - Unknown

there should be a help menu available with ? that will show all keyboard commands for that specific locality.

For example, the Applications Pane and the Pod Logs Pane would have different commands viewable in a modal with (?)

## Searchable panes - Medium

Every pane should have a local "/" command that will automatically search for relevant info

1. Case sensitive by default
2. Searches across multiple fields "smartly", such as Kind, Name, and Status
3. Stage 1 Substring search

## Pagination - Medium

Every list / search function from ArgoCD should have a pagination limit

# Nice to haves

## Login methods - Unknown - Nice to have

Logging in is restricted to ARGOCD_SERVER constant.

Some interesting ideas:

- env vars for SERVER_URL
- auto port forwarding / pass discovery
- passkey storage for passwords?
- SSO logins? For now no
- other methods?

## Grafana / Grafana Cloud - Unknownm - Nice to have

If this software were to go anywhere, it would be awesome to hook up
tracing, metrics, and logging to understand usage, limitations, catch bugs, and so on.

# Fun "unimportant" features

## Fuzzy Searching - Unknown - Not important

Fuzzy searching might be nice in certain areas

## Pinning / Marking - Unknown - Not important

It's an interesting idea to pin / mark certain applications like harpoon to fast switch maybe

Probably not useful.

# Todos

- CI/CD Pipeline
- Install script
- Update script

# Bug List

- first bug?
