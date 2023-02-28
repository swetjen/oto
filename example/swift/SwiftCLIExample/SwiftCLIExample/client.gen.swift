//  Code generated by oto; DO NOT EDIT.

import Foundation

class OtoClient {
	var endpoint: String
	init(withEndpoint url: String) {
		self.endpoint = url
	}
}


// GreeterService is a polite API for greeting people.
class GreeterService {
	var client: OtoClient
	init(withClient client: OtoClient) {
		self.client = client
	}

	// Greet prepares a lovely greeting.
	func greet(withRequest greetRequest: GreetRequest, completion: @escaping (_ response: GreetResponse?, _ error: Error?) -> ()) {
		let url = "\(self.client.endpoint)/GreeterService.Greet"
		var request = URLRequest(url: URL(string: url)!)
		request.httpMethod = "POST"
		request.addValue("application/json; charset=utf-8", forHTTPHeaderField: "Content-Type")
		request.addValue("application/json; charset=utf-8", forHTTPHeaderField: "Accept")
		var jsonData: Data
		do {
			jsonData = try JSONEncoder().encode(greetRequest)
		} catch let err {
			completion(nil, err)
			return
		}
		request.httpBody = jsonData
		let session = URLSession(configuration: URLSessionConfiguration.default)
		let task = session.dataTask(with: request) { (data, response, error) in
			if let err = error {
				completion(nil, err)
				return
			}
            if let httpResponse = response as? HTTPURLResponse {
                if (httpResponse.statusCode != 200) {
                    let err = OtoError("\(url): \(httpResponse.statusCode) status code")
                    completion(nil, err)
                    return
                }
            }
			var greetResponse: GreetResponse
			do {
				greetResponse = try JSONDecoder().decode(GreetResponse.self, from: data!)
			} catch let err {
				completion(nil, err)
				return
			}
            if let serviceErr = greetResponse.error {
                if (serviceErr != "") {
                    let err = OtoError(serviceErr)
                        completion(nil, err)
                        return
                }
            }
			completion(greetResponse, nil)
		}
		task.resume()
	}

}



// GreetRequest is the request object for GreeterService.Greet.
struct GreetRequest: Encodable, Decodable {

	// Name is the person to greet. It is required.
	var name: String?

}

// GreetResponse is the response object containing a person's greeting.
struct GreetResponse: Encodable, Decodable {

	// Greeting is a nice message welcoming somebody.
	var greeting: String?

	// Error is string explaining what went wrong. Empty if everything was fine.
	var error: String?

}


struct OtoError: LocalizedError
{
    var errorDescription: String? { return message }
    var failureReason: String? { return message }
    var recoverySuggestion: String? { return "" }
    var helpAnchor: String? { return "" }

    private var message : String

    init(_ description: String) {
        message = description
    }
}