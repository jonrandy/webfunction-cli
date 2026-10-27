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
 * Generated client for the Generics Convert Test webfunction package.
 * Source: https://example.test/package
 *
 * @phpstan-type FilterOfStringArgs = array{value: string}
 * @psalm-type FilterOfStringArgs = array{value: string}
 * @phpstan-type UserAttributes = array{id: string}
 * @psalm-type UserAttributes = array{id: string}
 * @phpstan-type ListOfUserAttributes = array{total: int|float, items: list<UserAttributes>}
 * @psalm-type ListOfUserAttributes = array{total: int|float, items: list<UserAttributes>}
 * @phpstan-type ListOfUserOrNullAttributes = array{total: int|float, items: list<UserAttributes|null>}
 * @psalm-type ListOfUserOrNullAttributes = array{total: int|float, items: list<UserAttributes|null>}
 * @phpstan-type PageOfUserAttributes = array{results: list<UserAttributes>}
 * @psalm-type PageOfUserAttributes = array{results: list<UserAttributes>}
 * @phpstan-type PageOfPageOfUserAttributes = array{results: list<PageOfUserAttributes>}
 * @psalm-type PageOfPageOfUserAttributes = array{results: list<PageOfUserAttributes>}
 * @phpstan-type PairOfNumberAttributes = array{first: int|float, second: int|float}
 * @psalm-type PairOfNumberAttributes = array{first: int|float, second: int|float}
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
     * @param array{filter?: FilterOfStringArgs} $args
     * @return ListOfUserAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listUsers(array $args = []): array
    {
        return $this->client->call('list-users', $args);
    }

    /**
     * @return ListOfUserOrNullAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listOptionalUsers(): array
    {
        return $this->client->call('list-optional-users');
    }

    /**
     * @return PageOfPageOfUserAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function nestedPages(): array
    {
        return $this->client->call('nested-pages');
    }

    /**
     * @return list<PairOfNumberAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function pairs(): array
    {
        return $this->client->call('pairs');
    }

}

