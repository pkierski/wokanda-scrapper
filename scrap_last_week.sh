#!/bin/sh

DIR=/root/sandbox

cd $DIR
$DIR/wokanda-scrapper -insecure-skip-verify
$DIR/wokanda-scrapper -insecure-skip-verify -relative-scrap-date=-1
$DIR/wokanda-scrapper -insecure-skip-verify -relative-scrap-date=-2
$DIR/wokanda-scrapper -insecure-skip-verify -relative-scrap-date=-3
$DIR/wokanda-scrapper -insecure-skip-verify -relative-scrap-date=-4
$DIR/wokanda-scrapper -insecure-skip-verify -relative-scrap-date=-5
$DIR/wokanda-scrapper -insecure-skip-verify -relative-scrap-date=-6
