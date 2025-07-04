<?php

$config = [
    [
        'name' => 'topic1',
        'partition' => 0,
        'worker' => 'worker.php',

    ],
    [
        'name' => 'topic2',
        'partition' => 0,
        'worker' => 'worker.php',

    ],
];


// Определяем C-функцию, объявленную в библиотеке
$ffi = FFI::cdef(
    'bool Consume(char* value);',
    './multiThreadConsumer.so'
);




// Вызываем функцию Add из библиотеки Go
$result = $ffi->Consume(json_encode($config));
echo $result;

