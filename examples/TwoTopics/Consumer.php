<?php
// Configuring the consumer
$config = [
    'kafka_connect' => 'localhost:9092',
    // disables/enables autocommit
    'commit_interval' => 0,
    // Prints debug information
    'debug_mode' => true,
    // For each partition, its own consumer is launched in its own goroutine.
    'workers' => [
        [
            'name' => 'topic1',
            'partition' => 0,
            'worker' => 'worker1.php',
            'group_id' => 'my-group1'
        ],
        [
            'name' => 'topic2',
            'partition' => 0,
            'worker' => 'worker2.php',
            'group_id' => 'my-group2'
        ],
    ]
];


// Определяем C-функцию, объявленную в библиотеке
$ffi = FFI::cdef(
    'bool MultiConsume(char* kafkaConnect, int commitInterval, bool debugMode,  char* jsonStr);',
    './../../bin/multiThreadConsumer.so'
);


// Вызываем функцию Add из библиотеки Go
$result = $ffi->MultiConsume(
    $config['kafka_connect'],
    (string)$config['commit_interval'],
    $config['debug_mode'],
    json_encode($config['workers'])
);
echo $result;

