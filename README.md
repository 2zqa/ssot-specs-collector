# SSOT specs collector

Collects information about hardware and sends it to a central server

> [!NOTE]
> This project is part of a suite of projects that work together. For all other related projects, see this search query: [`owner:2zqa topic:ssot`](https://github.com/search?q=owner%3A2zqa+topic%3Assot&type=repositories)

## Getting started

### Prerequisites

- Go (tested with 1.19)
- systemd-detect-virt
- [ssot-specs-server](https://github.com/2zqa/ssot-specs-server)

### Setup

1. Create a configuration file in `/etc/ssot-specs-collector/config`. For example:

    ```
    api-key 06a85GDe8pOgVTxm
    uuid d7ca6fb2-77af-40db-a0bc-33962df35bf5
    ```

2. Ensure that [ssot-specs-server](https://github.com/2zqa/ssot-specs-server) is running and reachable

### Installation and running

1. Clone and enter the repository: `git clone https://github.com/2zqa/ssot-specs-server.git && cd ssot-specs-server`
2. Run `go install`
3. Run `ssot-specs-collector`

> [!TIP]
> It is recommended to run the program as a service. For example, using systemd:
>
> ```
> [Unit]
> Description=SSOT specs collector
>
> [Service]
> ExecStart=/usr/bin/ssot-specs-collector
>
> [Install]
> WantedBy=multi-user.target
> ```

## License

SSOT specs collector is licensed under the [MIT](LICENSE) license.

## Acknowledgements

- [Voys](https://www.voys.nl/) for facilitating the internship where this project was developed
