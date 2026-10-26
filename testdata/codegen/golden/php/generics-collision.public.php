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
 * Generated client for the Generics Collision Test webfunction package.
 * Source: https://example.test/package
 *
 * @phpstan-type UserAttributes = array{id: string}
 * @psalm-type UserAttributes = array{id: string}
 * @phpstan-type ListOfUserAttributes = array{items: list<UserAttributes>}
 * @psalm-type ListOfUserAttributes = array{items: list<UserAttributes>}
 * @phpstan-type ListOfUserAttributes2 = array{old_field: bool}
 * @psalm-type ListOfUserAttributes2 = array{old_field: bool}
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
     * @return ListOfUserAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listUsers(): array
    {
        return $this->client->call('list-users');
    }

    /**
     * @return ListOfUserAttributes2
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getLegacy(): array
    {
        return $this->client->call('get-legacy');
    }

}

