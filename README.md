![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/jeffmay/starq/build.yaml)
[![CodeCov](https://codecov.io/gh/jeffmay/starq/graph/badge.svg?token=2aYYfHqAJW)](https://codecov.io/gh/jeffmay/starq)
![GitHub License](https://img.shields.io/github/license/jeffmay/starq)

# StarQ
(Pronounced "star-q")

A [`jq`](https://jqlang.github.io/jq/manual/) wrapper that uses a set of configured rules to transform an incoming YAML or JSON document into an outgoing YAML or JSON document.

## Contributing

[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)

This repo uses the [commitizen command-line tool](https://commitizen.github.io/cz-cli/) (`npx cz`) to write [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/#summary) messages. These commit messages are also run through the [devmoji](https://github.com/folke/devmoji?tab=readme-ov-file#sparkles-devmoji) editor / linter (`npx devmoji -e --lint`) to add gitmoji-style emoji (based on the commit type) and verify that the message follows the conventional format.

You should always initiate pull requests as drafts until you have addressed all of the bullet points on the pull request description and then, when you click "Ready to Review", it will tag the code owners to review the changes. You should not need to update the pull request description much after that -- other than to add screenshots, description of the bug or reason for the fix, etc -- as the pull request will use the conventional commits to generate the description.

Thank you for your support!
