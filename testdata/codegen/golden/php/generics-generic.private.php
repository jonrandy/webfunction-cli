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
 * Generated client for the Generics Codegen Test webfunction package.
 * Source: https://example.test/package
 *
 * @phpstan-type FilterOfStringArgs = array{value: string, operator?: string}
 * @psalm-type FilterOfStringArgs = array{value: string, operator?: string}
 * @phpstan-type UserAttributes = array{id: string, name: string, recovery_code: string}
 * @psalm-type UserAttributes = array{id: string, name: string, recovery_code: string}
 * @phpstan-type ListOfUserAttributes = array{items: list<UserAttributes>, total: int, cursor: string}
 * @psalm-type ListOfUserAttributes = array{items: list<UserAttributes>, total: int, cursor: string}
 * @phpstan-type ListOfUserOrNullAttributes = array{items: list<UserAttributes|null>, total: int, cursor: string}
 * @psalm-type ListOfUserOrNullAttributes = array{items: list<UserAttributes|null>, total: int, cursor: string}
 * @phpstan-type PageOfUserAttributes = array{results: list<UserAttributes>, has_more: bool}
 * @psalm-type PageOfUserAttributes = array{results: list<UserAttributes>, has_more: bool}
 * @phpstan-type PageOfPageOfUserAttributes = array{results: list<PageOfUserAttributes>, has_more: bool}
 * @psalm-type PageOfPageOfUserAttributes = array{results: list<PageOfUserAttributes>, has_more: bool}
 * @phpstan-type PageOfArrayOfUserAttributes = array{results: list<list<UserAttributes>>, has_more: bool}
 * @psalm-type PageOfArrayOfUserAttributes = array{results: list<list<UserAttributes>>, has_more: bool}
 * @phpstan-type PairOfNumberAttributes = array{first: int|float, second: int|float}
 * @psalm-type PairOfNumberAttributes = array{first: int|float, second: int|float}
 * @phpstan-type BothOfStringArgs = array{input?: string}
 * @psalm-type BothOfStringArgs = array{input?: string}
 * @phpstan-type BothOfNumberAttributes = array{output: int|float}
 * @psalm-type BothOfNumberAttributes = array{output: int|float}
 * @phpstan-type HolderAttributes = array{users: ListOfUserAttributes, tags: list<string>}
 * @psalm-type HolderAttributes = array{users: ListOfUserAttributes, tags: list<string>}
 * @phpstan-type PairOfStringAttributes = array{first: string, second: string}
 * @psalm-type PairOfStringAttributes = array{first: string, second: string}
 * @phpstan-type WrapperOfStringAttributes = array{items: list<string>, pair: PairOfStringAttributes}
 * @psalm-type WrapperOfStringAttributes = array{items: list<string>, pair: PairOfStringAttributes}
 * @phpstan-type AuditAttributes = array{at: int|float, note: string}
 * @psalm-type AuditAttributes = array{at: int|float, note: string}
 * @phpstan-type PageOfAuditAttributes = array{results: list<AuditAttributes>, has_more: bool}
 * @psalm-type PageOfAuditAttributes = array{results: list<AuditAttributes>, has_more: bool}
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
     * @param array{filter?: FilterOfStringArgs, debug_trace?: bool} $args
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
    public function listUsersOrNull(): array
    {
        return $this->client->call('list-users-or-null');
    }

    /**
     * @return PageOfPageOfUserAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function pageOfPages(): array
    {
        return $this->client->call('page-of-pages');
    }

    /**
     * @return PageOfArrayOfUserAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function pageOfArrays(): array
    {
        return $this->client->call('page-of-arrays');
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

    /**
     * @param array{payload: BothOfStringArgs} $args
     * @return BothOfNumberAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function convert(array $args): array
    {
        return $this->client->call('convert', $args);
    }

    /**
     * @return HolderAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getHolder(): array
    {
        return $this->client->call('get-holder');
    }

    /**
     * @return WrapperOfStringAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getWrapper(): array
    {
        return $this->client->call('get-wrapper');
    }

    /**
     * @return PageOfAuditAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listAudit(): array
    {
        return $this->client->call('list-audit');
    }

}

