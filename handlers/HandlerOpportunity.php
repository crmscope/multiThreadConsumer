<?php

/**
 * Пример запроса из Битрикса (Отформатированный).
 * {
 *    "id": 4781,
 *    "title": "",
 *    "companyId": 0,
 *    "stageId": "C2:EXECUTING",
 *    "sourceId": "",
 *    "comments": null,
 *    "opportunity": 0,
 *    "account": {
 *        "id": 0,
 *        "title": null,
 *        "companyInn": null,
 *        "companyKpp": null,
 *        "companyType": null,
 *        "industry": null,
 *        "comments": null
 *    }
 * }
 *
 */
$request = base64_decode(readline(""));
$param =  json_decode($request,1);

// Пишем параметры в какой ни будь лог
file_put_contents('/var/log/crm/TEMPORARY_BITRIX_LOG.txt', date('Y-m-d H:i:s') . " " . print_r($param,1) . PHP_EOL, FILE_APPEND);
exit();
