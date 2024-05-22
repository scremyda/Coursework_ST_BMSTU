# Coursework_ST_BMSTU

## Канальный уровень
Данный уровень эмитирует взаимодействие с удаленным сетевым узлом через канал с помехами. Для этого используется один из четырех вариантов кодирования передаваемых сообщений. Полученный от транспортного уровня json сегмента кодируется в битовый формат соответствующим кодом, вносятся ошибки и затем декодируется для последующей отправки обратно на транспортный уровень. Ошибка вносится в 1 случайный бит сегмента.

Сервис должен вносить ошибку с вероятностью P в передаваемую информацию и терять сообщения с вероятностью R. При декодировании либо пакет с ошибкой теряется, либо передается обратно на транспортный после исправления.

## Вариант 11
Отправка текста. Передаваемую информацию защитить передаваемую информацию [7,4]-кодом Хэмминга. Длина сегмента (X) 150 байт, период сборки сегментов (N) 2 секунд, вероятность ошибки (P) 8%, вероятность потери кадра (R) 1%.
