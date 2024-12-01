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
    response = requests.get("http://localhost:8080/login", json={'username': 'test', 'password': 'test'})
    assert response.status_code == 200

def test_upload():
    #primero, necesitamos el token
    response = requests.get("http://localhost:8080/login", json={'username':'test', 'password':'test'})
    token = response.json()['token']
    #ahora, subimos el archivo
    response = requests.post("http://localhost:8080/test/archivo", headers={'Authorization': token}, json={'doc_content': 'test'})
    assert response.status_code == 201
    response = requests.post("http://localhost:8080/test/archivo", headers={'Authorization': token}, json={'doc_content': 'test'})
    assert response.status_code == 409

def test_get():
    #primero, necesitamos el token
    response = requests.get("http://localhost:8080/login", json={'username':'test', 'password':'test'})
    token = response.json()['token']
    response = requests.get("http://localhost:8080/test/archivo", headers={'Authorization': token})
    assert response.status_code == 200
    response = requests.get("http://localhost:8080/test/archivo2", headers={'Authorization': token})
    assert response.status_code == 404

def test_list():
    #primero, necesitamos el token
    response = requests.get("http://localhost:8080/login", json={'username':'test', 'password':'test'})
    token = response.json()['token']
    requests.post("http://localhost:8080/test/archivo2", headers={'Authorization': token}, json={'doc_content': 'test'})
    response = requests.get("http://localhost:8080/test/_all_docs", headers={'Authorization': token})
    assert response.status_code == 200

def test_update():
    response = requests.get("http://localhost:8080/login", json={'username':'test', 'password':'test'})
    token = response.json()['token']
    response = requests.put("http://localhost:8080/test/archivo", headers={'Authorization': token}, json={'doc_content': 'test'})
    assert response.status_code == 200

def test_delete():
    response = requests.get("http://localhost:8080/login", json={'username':'test', 'password':'test'})
    token = response.json()['token']
    response = requests.delete("http://localhost:8080/test/archivo", headers={'Authorization': token})
    assert response.status_code == 204
    response = requests.delete("http://localhost:8080/test/archivo", headers={'Authorization': token})
    assert response.status_code == 404
