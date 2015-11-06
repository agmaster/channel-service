#!/usr/bin/env ruby
## encoding: utf-8
#
#

require 'elasticsearch'

while true
  client = Elasticsearch::Client.new(log: true)
  client.transport.reload_connections!
  client.cluster.health
  client = nil
  sleep 1
end
