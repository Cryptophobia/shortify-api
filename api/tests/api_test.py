import requests
import json

# Add one url
data = {'longurl': 'https://github.com/teamhephy/workflow/'}
r = requests.post('http://localhost:5000/Create', data=json.dumps(data))
print("Status code received:" + str(r.status_code) + " response body: " + r.text)

# Add another url
data = {'longurl': 'https://teamhephy.com'}
r = requests.post('http://localhost:5000/Create', data=json.dumps(data))
print("Status code received:" + str(r.status_code) + " response body: " + r.text)

# Find the recently added shortURL from the response
json_data = r.json()
targetURL = str.format("http://localhost:5000/{}", json_data["ShortURL"])
r = requests.get(targetURL)
print("Status code from get request:" + str(r.status_code) + " response body: " + r.text)

# Add another invalid url
data = {'longurl': 'https//teamhephy.com'}
r = requests.post('http://localhost:5000/Create', data=json.dumps(data))
print("Status code received:" + str(r.status_code) + " response body: " + r.text)

# Try a shortURL that is not long enough
data = {'ShortURL': '8sa3efd'}
targetURL = str.format("http://localhost:5000/{}", data["ShortURL"])
r = requests.get(targetURL)
print("Status code received:" + str(r.status_code) + " response body: " + r.text)
