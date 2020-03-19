require 'net/http'
require 'json'
require 'awesome_print'

# Client implements a ruby interface to the Population API service.
class Client
  BASE_URI = 'https://server-7y3morjijq-uw.a.run.app/zip/'

  # Retrive population metadata for a given zip code.
  def zip(zip_code)
    uri = URI(BASE_URI + zip_code)
    JSON.parse Net::HTTP.get(uri)
  end
end

c = Client.new
ap c.zip(ARGV[0])
