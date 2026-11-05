ABC Financial mock back ends - Lab 04 (Creating an API Product)
=================================================================

One small program that starts two separate back ends, one for each Lab 04 API,
so the CustomerInformation API, the Leasing API and the CustomerLeasing API
Product all return real data. The API Product shows clearly that one product
can reach two different back ends.

Customer Info back end               Leasing back end
http://localhost:8091                http://localhost:8092
  GET    /customers                    GET    /assets
  POST   /customers                    POST   /assets
  GET    /customers/{customerId}       GET    /assets/{assetId}
  PUT    /customers/{customerId}       PUT    /assets/{assetId}
  DELETE /customers/{customerId}       DELETE /assets/{assetId}

Sample data: customers 1001, 1002, 1003 and assets (leases) 5001, 5002, 5003.
Data is kept in memory and resets every time the server is restarted.


1. Start the back end
---------------------
Pick the file for your laptop from the bin folder. No installation is needed.

  macOS (Apple silicon M1/M2/M3/M4):  ./abc-backend-mac-arm64
  macOS (Intel):                      ./abc-backend-mac-intel
  Windows:                            abc-backend-windows.exe
  Linux:                              ./abc-backend-linux

The program starts both back ends and prints:
  Customer Info back end listening on http://localhost:8091
  Leasing back end listening on http://localhost:8092

Every request is logged with the name of the back end that received it, for
example [Customer Info] GET /customers/1001 or [Leasing] GET /assets/5001.
Leave the window open while you do the lab. Stop it with Ctrl+C.

If port 8091 or 8092 is already in use, choose other ports, for example:
  ./abc-backend-mac-arm64 -customer-port 9091 -leasing-port 9092

If you have Go 1.22 or later installed, you can also run it from source:
  go run main.go

macOS may block the file the first time ("cannot be opened because the
developer cannot be verified"). Run this once in the bin folder, then start
it again:
  xattr -d com.apple.quarantine abc-backend-mac-*
On Windows, if SmartScreen appears, choose More info, then Run anyway.

Check it works: open http://localhost:8091/customers (three customers) and
http://localhost:8092/assets (three leased assets) in the browser.


2. Use it in Lab 04
-------------------
Each API points to its own back end. Enter these as the Endpoint when you
create the APIs (production and sandbox):

  CustomerInformation API   http://localhost:8091
  Leasing API               http://localhost:8092

The gateway adds the resource path to the endpoint, so for example:
  GET /customer-info/1.0.0/customers/1001  -> Customer Info  GET /customers/1001
  GET /leasing/1.0.0/assets/5001           -> Leasing        GET /assets/5001
  GET /customer-leasing/1.0.0/customers    -> Customer Info  GET /customers  (API Product)
  GET /customer-leasing/1.0.0/assets/5001  -> Leasing        GET /assets/5001 (API Product)

When API Manager runs on a different machine from the back ends, replace
localhost with the IP address or host name of the machine running them.


3. Try it through the gateway
-----------------------------
After subscribing to the CustomerLeasing product and generating a token in the
Developer Portal:

  curl -k -H "Authorization: Bearer <ACCESS_TOKEN>" \
    https://<APIM_HOST>:<GW_HTTPS_PORT>/customer-leasing/1.0.0/customers

  curl -k -H "Authorization: Bearer <ACCESS_TOKEN>" \
    https://<APIM_HOST>:<GW_HTTPS_PORT>/customer-leasing/1.0.0/customers/1001

  curl -k -H "Authorization: Bearer <ACCESS_TOKEN>" \
    https://<APIM_HOST>:<GW_HTTPS_PORT>/customer-leasing/1.0.0/assets/5001

  curl -k -X POST -H "Authorization: Bearer <ACCESS_TOKEN>" \
    -H "Content-Type: application/json" \
    -d '{"customerId":"1001","type":"VEHICLE","description":"Mazda CX-5","leaseAmount":2500000,"monthlyRental":75000,"termMonths":36}' \
    https://<APIM_HOST>:<GW_HTTPS_PORT>/customer-leasing/1.0.0/assets

Resources that are not part of the product (for example GET /assets or
DELETE /customers/{customerId}) are rejected by the gateway and never reach the
back end - watch the back end window to see this. The product calls that are
allowed show up under both [Customer Info] and [Leasing].
