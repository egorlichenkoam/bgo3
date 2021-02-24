INSERT INTO clients
    (login, "password", full_name, passport, birthday, status)
VALUES ('pupok', '123456', 'pupkin pupik pupkovich', '1234567890', '2011-11-11', 'INACTIVE');
INSERT INTO clients
    (login, "password", full_name, passport, birthday, status)
VALUES ('shlupka', '654321', 'shlupkin shlupka shlupkovich', '0987654321', '2011-11-11', 'ACTIVE');
INSERT INTO clients
    (login, "password", full_name, passport, birthday, status)
VALUES ('shlupka2', '654321', 'shlupkina shlupka shlupkovich', '0987654321', '2011-11-11', 'ACTIVE');

INSERT INTO cards
    ("number", balance, issuer, holder, owner_id, status)
VALUES ('1234', 1000000, 'VISA', 'PUP', 1, 'ACTIVE');
INSERT INTO cards
    ("number", balance, issuer, holder, owner_id, status)
VALUES ('4321', 333333, 'MasterCard', 'SHLUP', 1, 'ACTIVE');
INSERT INTO cards
    ("number", balance, issuer, holder, owner_id, status)
VALUES ('5678', 2222222, 'VISA', 'PUP', 2, 'ACTIVE');
INSERT INTO cards
    ("number", balance, issuer, holder, owner_id, status)
VALUES ('8765', 44444444, 'MIR', 'SHLUPKA', 3, 'ACTIVE');

INSERT INTO icons
    (url)
VALUES ('1');
INSERT INTO icons
    (url)
VALUES ('2');
INSERT INTO icons
    (url)
VALUES ('3');
INSERT INTO icons
    (url)
VALUES ('4');

INSERT INTO transactions
    (card_id, amount, tx_type, "comments", mcc, icon_id, status)
VALUES (1, 5000000, 'TO', 'Пополнение через Альфа-Банк', '1234', 1, 'EXECUTED');
INSERT INTO transactions
    (card_id, amount, tx_type, "comments", mcc, icon_id, status)
VALUES (1, 100000, 'FROM', 'Продукты', '5678', 2, 'EXECUTED');
INSERT INTO transactions
    (card_id, amount, tx_type, "comments", mcc, icon_id, status)
VALUES (1, 100000, 'FROM', 'Пополнение мобильного телефона', '9012', 3, 'EXECUTED');
INSERT INTO transactions
    (card_id, amount, tx_type, "comments", mcc, icon_id, status)
VALUES (1, 100000, 'FROM', 'Перевод', '3456', 4, 'EXECUTED');

COMMIT;