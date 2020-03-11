## Функциональность
Проксирование http

Проксирование https

Отправка сохраненного запроса

## Как проверить работу
(Использовала mozilla firefox)

Добавить сертификат в rootCert.cert в браузер

Добавить 127.0.0.1:5000 как адрес прокси сервера

Выполнить команды

sudo docker build -t alex https://github.com/saskamegaprogrammist/proxyServer.git

sudo docker run -p 5000:5000 -p 5001:5001 --name alex -t alex

## Функциональность repeater
localhost:5001/requests - выдает информацию о последних 10 запросах

localhost:5001/requests/{id} - повторяет запрос по id запроса


