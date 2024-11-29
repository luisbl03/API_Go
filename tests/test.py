import pytest 
import requests

def test_version():
    response = requests.get("http://localhost:8080/version")
    assert response.status_code == 200
    assert response.json() == {"version": "1.0.0"}

def test_signup():
    response = requests.post("http://localhost:8080/signup", json= {'username': 'test', 'password': 'test'})
    print(response.json())
    assert response.status_code == 201
    response = requests.post("http://localhost:8080/signup", json={'username': 'test', 'password': 'test'})
    print(response.json())
    assert response.status_code == 409

def test_login():
    response = requests.post("http://localhost:8080/login", json= {'username': 'test', 'password': 'test'})
    print(response.json())
    assert response.status_code == 200

def test_upload():
    response = requests.get("http://localhost:8080/login", json= {'username': 'test', 'password': 'test'})
    token = response.json()['token']
    response = requests.post("http://localhost:8080/test/archivo", headers={'Authorization':token}, json={'doc_content': 'test'})
    print(response.json())
    assert response.status_code == 201
    response = requests.post("http://localhost:8080/test/archivo", headers={'Authorization':token}, json={'doc_content': 'test'})
    print(response.json())
    assert response.status_code == 409

def test_get():
    response = requests.get("http://localhost:8080/test/archivo", headers={'Authorization':token})
    print(response.json())
    assert response.status_code == 200
    assert response.json() == {'doc_content': 'test'}

def test_list():
    response = requests.get("http://localhost:8080/test/_all_docs", headers={'Authorization':token})
    print(response.json())
    assert response.status_code == 200

def test_update():
    response = requests.put("http://localhost:8080/test/archivo", headers={'Authorization':token}, json={'doc_content': 'test2'})
    print(response.json())
    assert response.status_code == 200
    response = requests.get("http://localhost:8080/test/archivo", headers={'Authorization':token})
    print(response.json())
    assert response.status_code == 200
    assert response.json() == {'doc_content': 'test2'}

def test_delete():
    response = requests.delete("http://localhost:8080/test/archivo", headers={'Authorization':token})
    print(response.json())
    assert response.status_code == 204
    response = requests.delete("http://localhost:8080/test/archivo", headers={'Authorization':token})
    print(response.json())
    assert response.status_code == 404