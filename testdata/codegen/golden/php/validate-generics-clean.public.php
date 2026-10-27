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
 * Generated client for the Generics Clean API webfunction package.
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
 * @phpstan-type HolderAttributes = array{inner: ListOfUserAttributes}
 * @psalm-type HolderAttributes = array{inner: ListOfUserAttributes}
 * @phpstan-type ListOfListOfUserAttributes = array{total: int|float, items: list<ListOfUserAttributes>}
 * @psalm-type ListOfListOfUserAttributes = array{total: int|float, items: list<ListOfUserAttributes>}
 * @phpstan-type NestedHolderAttributes = array{deep: ListOfListOfUserAttributes}
 * @psalm-type NestedHolderAttributes = array{deep: ListOfListOfUserAttributes}
 * @phpstan-type CommentAttributes = array{message: string, replies: list<CommentAttributes>}
 * @psalm-type CommentAttributes = array{message: string, replies: list<CommentAttributes>}
 * @phpstan-type ThreadAttributes = array{comments: ListOfThreadAttributes}
 * @psalm-type ThreadAttributes = array{comments: ListOfThreadAttributes}
 * @phpstan-type ListOfThreadAttributes = array{total: int|float, items: list<ThreadAttributes>}
 * @psalm-type ListOfThreadAttributes = array{total: int|float, items: list<ThreadAttributes>}
 * @phpstan-type PairOfStringAttributes = array{first: string, second: string}
 * @psalm-type PairOfStringAttributes = array{first: string, second: string}
 * @phpstan-type WrapperOfStringAttributes = array{pairs: list<PairOfStringAttributes>, items: list<string>}
 * @psalm-type WrapperOfStringAttributes = array{pairs: list<PairOfStringAttributes>, items: list<string>}
 * @phpstan-type PageOfArrayOfUserAttributes = array{results: list<list<UserAttributes>>}
 * @psalm-type PageOfArrayOfUserAttributes = array{results: list<list<UserAttributes>>}
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

    /**
     * @return list<HolderAttributes|NestedHolderAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getHolders(): array
    {
        return $this->client->call('get-holders');
    }

    /**
     * @return list<CommentAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getComments(): array
    {
        return $this->client->call('get-comments');
    }

    /**
     * @return ListOfThreadAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getThreads(): array
    {
        return $this->client->call('get-threads');
    }

    /**
     * @return WrapperOfStringAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function wrapped(): array
    {
        return $this->client->call('wrapped');
    }

    /**
     * @return PageOfArrayOfUserAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function pageOfUsers(): array
    {
        return $this->client->call('page-of-users');
    }

    /**
     * @return list<UserAttributes>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function users(): array
    {
        return $this->client->call('users');
    }

}

