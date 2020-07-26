import os
import json
import requests

if __name__ == "__main__":
    port = os.getenv('ADDRESSBOOK_PORT', '8081')
    address = 'http://localhost:' + port

    username = 'user'
    password = 'password'

    print('Creating a user')
    payload = {'username': username,
               'password': password, 'email': 'some@email.com'}
    print('POSTing to /api/user with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/user', json=payload)
    print('Received reply: ' + str(response.content))
    print('')

    print('Creating a token')
    payload = {'username': username, 'password': password}
    print('POSTing to /api/user/token with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/user/token', json=payload)
    print('Received reply: ' + str(response.content))
    token = "Bearer " + json.loads(response.content)["token"]
    print("Extracted token and using it in Authorization header for all subsequent requests: " + token)
    headers = {'Authorization': token}
    print('')

    print('Creating a contact')
    payload = {'name': 'name', 'surname': 'surname', 'email': 'valid@mail.com'}
    print('POSTing to /api/contact with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Creating another contact')
    payload = {'name': 'name2', 'surname': 'surname2',
               'email': 'valid2@mail.com'}
    print('POSTing to /api/contact with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    contactID1 = json.loads(response.content)["id"]
    print('')

    print('Creating another contact')
    payload = {'name': 'name3', 'surname': 'surname3',
               'email': 'valid3@mail.com'}
    print('POSTing to /api/contact with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    contactID2 = json.loads(response.content)["id"]
    print('')

    print('Creating another contact')
    payload = {'name': 'name3', 'surname': 'surname3',
               'email': 'valid3@mail.com'}
    print('POSTing to /api/contact with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    contactID3 = json.loads(response.content)["id"]
    print('')

    print('Fetching all contacts')
    print('GETing /api/contact')
    response = requests.get(address + '/api/contact', headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Fetching single contact with ID of ' + str(contactID3))
    print('GETing /api/contact/' + str(contactID3))
    response = requests.get(address + '/api/contact/' +
                            str(contactID3), headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Deleting a contact with ID of ' + str(contactID3))
    print('DELETEing /api/contact/' + str(contactID3))
    response = requests.delete(address + '/api/contact/' +
                               str(contactID3), headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Fetching all contacts')
    print('GETing /api/contact')
    response = requests.get(address + '/api/contact', headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Creating a contact-list')
    payload = {'name': 'name'}
    print('POSTing to /api/contact-list with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact-list',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    contactListID1 = json.loads(response.content)["id"]
    print('')

    print('Creating a contact-list')
    payload = {'name': 'name2'}
    print('POSTing to /api/contact-list with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact-list',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    contactListID2 = json.loads(response.content)["id"]
    print('')

    print('Creating another contact-list')
    payload = {'name': 'name3'}
    print('POSTing to /api/contact-list with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact-list',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    contactListID3 = json.loads(response.content)["id"]
    print('')

    print('Fetching all contact-lists')
    print('GETing /api/contact-list')
    response = requests.get(address + '/api/contact-list', headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Fetching single contact-list with ID of ' + str(contactListID3))
    print('GETing /api/contact-list/' + str(contactListID3))
    response = requests.get(address + '/api/contact-list/' +
                            str(contactListID3), headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Deleting a contact-list with ID of ' + str(contactListID3))
    print('DELETEing /api/contact-list/' + str(contactListID3))
    response = requests.delete(address + '/api/contact-list/' +
                               str(contactListID3), headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Fetching all contact-lists')
    print('GETing /api/contact-list')
    response = requests.get(address + '/api/contact-list', headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Creating another contact-list')
    payload = {'name': 'differentName'}
    print('POSTing to /api/contact-list with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact-list',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    contactListID4 = json.loads(response.content)["id"]
    print('')

    print('Conducting search on contact-lists')
    payload = {'term': 'name'}
    print('POSTing to /api/contact-list/search with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact-list/search',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Conducting another search on contact-lists')
    payload = {'term': 'different'}
    print('POSTing to /api/contact-list/search with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact-list/search',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Adding contact ID ' + str(contactID1) +
          ' to contact-list ID ' + str(contactListID1))
    payload = {'id': contactID1}
    print('POSTing to /api/contact-list/' + str(contactListID1) +
          '/contact with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact-list/'+str(contactListID1)+'/contact',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Adding contact ID ' + str(contactID2) +
          ' to contact-list ID ' + str(contactListID1))
    payload = {'id': contactID2}
    print('POSTing to /api/contact-list/' + str(contactListID1) +
          '/contact with payload: ' + json.dumps(payload))
    response = requests.post(address + '/api/contact-list/'+str(contactListID1)+'/contact',
                             json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Listing contacts of contact-list ID ' + str(contactListID1))
    print('GETing to /api/contact-list/' + str(contactListID1) + '/contact')
    response = requests.get(address + '/api/contact-list/' +
                            str(contactListID1) + '/contact', headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Deleting contact ID ' + str(contactID2) +
          ' from contact-list ID ' + str(contactListID1))
    payload = {'id': contactID2}
    print('DELETEing to /api/contact-list/' + str(contactListID1) +
          '/contact with payload: ' + json.dumps(payload))
    response = requests.delete(address + '/api/contact-list/'+str(contactListID1)+'/contact',
                               json=payload, headers=headers)
    print('Received reply: ' + str(response.content))
    print('')

    print('Listing contacts of contact-list ID ' + str(contactListID1))
    print('GETing to /api/contact-list/' + str(contactListID1) + '/contact')
    response = requests.get(address + '/api/contact-list/' +
                            str(contactListID1) + '/contact', headers=headers)
    print('Received reply: ' + str(response.content))
