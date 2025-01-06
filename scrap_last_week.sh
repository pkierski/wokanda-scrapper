#!/bin/sh

DIR=/root/sandbox

cd $DIR
$DIR/scrapper -insecure-skip-verify
$DIR/scrapper -insecure-skip-verify -relative-scrap-date=-1
$DIR/scrapper -insecure-skip-verify -relative-scrap-date=-2
$DIR/scrapper -insecure-skip-verify -relative-scrap-date=-3
$DIR/scrapper -insecure-skip-verify -relative-scrap-date=-4
$DIR/scrapper -insecure-skip-verify -relative-scrap-date=-5
$DIR/scrapper -insecure-skip-verify -relative-scrap-date=-6
