<?php

$response = [
    'error'=> 200,
    'errorMessage' => $request
];
echo json_encode(print_r($param,1), JSON_UNESCAPED_UNICODE);

