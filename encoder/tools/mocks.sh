#!/bin/bash

mockery -case=underscore -all -dir=pkg -output=pkg/mocks -outpkg=encoder_mocks
