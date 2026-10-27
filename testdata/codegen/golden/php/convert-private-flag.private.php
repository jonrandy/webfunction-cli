<?php

declare(strict_types=1);

namespace WebFunctionClient;

use WebFunction\Page;
use WebFunction\Package;
use WebFunction\Pipeline;
use WebFunction\BadRequestError;
use WebFunction\UnexpectedStatusCodeError;
use WebFunction\JsonParseError;

// IMPORTANT - composer require webfunction-protocol/webfunction

/**
 * Generated client for the Private Flag Convert Test webfunction package.
 * Source: https://example.test/package
 *
 * @phpstan-type FilterArgs = array{status?: string, hidden_override?: string}
 * @psalm-type FilterArgs = array{status?: string, hidden_override?: string}
 * @phpstan-type PersonAttributes = array{name: string, audit_ref: string}
 * @psalm-type PersonAttributes = array{name: string, audit_ref: string}
 */
final class Client
{
    private \WebFunction\Client $client;

    public function __construct(?string $bearerAuth = null, ?string $version = null)
    {
        $this->client = \WebFunction\Client::fromPackageEndpoint('https://example.test/package', bearerAuth: $bearerAuth, version: $version);
    }

    /**
     * Escape hatch for calling an endpoint this client doesn't have a typed method for.
     */
    public function call(string $name, array $args = []): mixed
    {
        return $this->client->call($name, $args);
    }

    public function getPackage(): ?Package
    {
        return $this->client->getPackage();
    }

    public function setBearerAuth(?string $bearerAuth): void
    {
        $this->client->setBearerAuth($bearerAuth);
    }

    public function setVersion(?string $version): void
    {
        $this->client->setVersion($version);
    }

    public function setPipeline(?Pipeline $pipeline): void
    {
        $this->client->setPipeline($pipeline);
    }

    /**
     * @param array{sync_secret?: string} $args
     * @return bool
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function internalSync(array $args = []): bool
    {
        return $this->client->call('internal-sync', $args);
    }

    /**
     * @param array{id: string, debug_trace?: bool, filter?: FilterArgs} $args
     * @return PersonAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getPerson(array $args): array
    {
        return $this->client->call('get-person', $args);
    }

    /**
     * @return list<array{text: string, internal_flag: bool}>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listNotes(): array
    {
        return $this->client->call('list-notes');
    }

}

