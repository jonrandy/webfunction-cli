# Generated client for the Generics Convert Test webfunction package.
# Source: https://example.test/package
#
# Wraps the real web_function gem (github.com/robinclart/web_function) -
# every method here calls straight through to the real dynamic client.
# frozen_string_literal: true

require "web_function"

module WebFunctionClient
  class UserAttributes
    def initialize(raw)
      @raw = raw
    end

    def id
      @raw["id"]
    end

  end

  class ListOfUserAttributes
    def initialize(raw)
      @raw = raw
    end

    def total
      @raw["total"]
    end

    def items
      (@raw["items"] || []).map { |item| UserAttributes.new(item) }
    end

  end

  class ListOfUserOrNullAttributes
    def initialize(raw)
      @raw = raw
    end

    def total
      @raw["total"]
    end

    def items
      (@raw["items"] || []).map { |item| UserAttributes.new(item) }
    end

  end

  class PageOfUserAttributes
    def initialize(raw)
      @raw = raw
    end

    def results
      (@raw["results"] || []).map { |item| UserAttributes.new(item) }
    end

  end

  class PageOfPageOfUserAttributes
    def initialize(raw)
      @raw = raw
    end

    def results
      (@raw["results"] || []).map { |item| PageOfUserAttributes.new(item) }
    end

  end

  class PairOfNumberAttributes
    def initialize(raw)
      @raw = raw
    end

    def first
      @raw["first"]
    end

    def second
      @raw["second"]
    end

  end

  # Sentinel distinguishing "caller omitted this optional argument" from
  # "caller explicitly passed nil" - forwarding nil either way would send
  # the real client (and the real server) an explicit JSON null for an
  # argument that was never actually given.
  Undefined = ::Object.new.freeze

  class Client
    def initialize(real_client)
      @real_client = real_client
    end

    # @param url [String]
    def self.from_package_endpoint(url, bearer_auth: nil, version: nil, pipelined: false)
      new(::WebFunction::Client.from_package_endpoint(url, bearer_auth: bearer_auth, version: version, pipelined: pipelined))
    end

    def list_users(filter: Undefined)
      args = {}
      args[:filter] = filter unless filter.equal?(Undefined)
      (v = @real_client.list_users(**args)) && ListOfUserAttributes.new(v)
    end

    def list_optional_users
      (v = @real_client.list_optional_users) && ListOfUserOrNullAttributes.new(v)
    end

    def nested_pages
      (v = @real_client.nested_pages) && PageOfPageOfUserAttributes.new(v)
    end

    def pairs
      (@real_client.pairs || []).map { |item| PairOfNumberAttributes.new(item) }
    end

    # Escape hatch for calling an endpoint this file doesn't declare a
    # typed method for.
    def call(endpoint_name, args = {})
      @real_client.call(endpoint_name, args)
    end

    def package
      @real_client.package
    end

    def bearer_auth=(value)
      @real_client.bearer_auth = value
    end

    def version=(value)
      @real_client.version = value
    end

    def pipeline=(value)
      @real_client.pipeline = value
    end

  end
end
