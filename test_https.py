import pytest 
import requests

cert = ('cmd/broker/certs/fullchain.pem', 'cmd/broker/certs/privkey.pem')

def test_version():
    response = requests.get("https://lbroker.duckdns.org:5000/version", cert=cert)
    assert response.status_code == 200
    assert response.json() == {"version": "1.0.0"}

def test_signup():
    response = requests.post("https://lbroker.duckdns.org:5000/signup", json={"username": "test", "password": "test"}, cert=cert)
    assert response.status_code == 201
    response = requests.post("https://lbroker.duckdns.org:5000/signup", json={"username": "test", "password": "test"}, cert=cert)
    assert response.status_code == 409

def test_login():
    response = requests.get("https://lbroker.duckdns.org:5000/login", json={"username": "test", "password": "test"}, cert=cert)
    assert response.status_code == 200
    response = requests.get("https://lbroker.duckdns.org:5000/login", json={"username": "test", "password": "test2"}, cert=cert)
    assert response.status_code == 401

def test_uploadJson():
    response = requests.get("https://lbroker.duckdns.org:5000/login", json={"username": "test", "password": "test"}, cert=cert)
    token = response.json()["token"]
    headerList = {
                "Authorization": token,
                "Content-Type": "application/json" }
    response = requests.post("https://lbroker.duckdns.org:5000/test/archivo", json={"doc_content":"test"}, cert=cert, headers=headerList)
    assert response.status_code == 201
    response = requests.post("https://lbroker.duckdns.org:5000/test/archivo", json={"doc_content":"test"}, cert=cert, headers=headerList)
    assert response.status_code == 409

def test_getJson():
    response = requests.get("https://lbroker.duckdns.org:5000/login", json={"username": "test", "password": "test"}, cert=cert)
    token = response.json()["token"]
    headerList = {
                "Authorization": token,
                "Content-Type": "application/json" }
    response = requests.get("https://lbroker.duckdns.org:5000/test/archivo", cert=cert, headers=headerList)
    assert response.status_code == 200
    response = requests.get("https://lbroker.duckdns.org:5000/test/archivo2", cert=cert, headers=headerList)
    assert response.status_code == 404

def test_updateJson():
    response = requests.get("https://lbroker.duckdns.org:5000/login", json={"username": "test", "password": "test"}, cert=cert)
    token = response.json()["token"]
    headerList = {
                "Authorization": token,
                "Content-Type": "application/json" }
    response = requests.put("https://lbroker.duckdns.org:5000/test/archivo", json={"doc_content":"test2"}, cert=cert, headers=headerList)
    assert response.status_code == 200
    response = requests.put("https://lbroker.duckdns.org:5000/test/archivo2", json={"doc_content":"test2"}, cert=cert, headers=headerList)
    assert response.status_code == 404

def test_list():
    response = requests.get("https://lbroker.duckdns.org:5000/login", json={"username": "test", "password": "test"}, cert=cert)
    token = response.json()["token"]
    headerList = {
                "Authorization": token,
                "Content-Type": "application/json" }
    response = requests.get("https://lbroker.duckdns.org:5000/test/_all_docs", cert=cert, headers=headerList)
    assert response.status_code == 200

def test_deleteJson():
    response = requests.get("https://lbroker.duckdns.org:5000/login", json={"username": "test", "password": "test"}, cert=cert)
    token = response.json()["token"]
    headerList = {
                "Authorization": token,
                "Content-Type": "application/json" }
    response = requests.delete("https://lbroker.duckdns.org:5000/test/archivo", cert=cert, headers=headerList)
    assert response.status_code == 204
    response = requests.delete("https://lbroker.duckdns.org:5000/test/archivo2", cert=cert, headers=headerList)
    assert response.status_code == 404