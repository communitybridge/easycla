#!/usr/bin/env python

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
import logging
import os
import sys
from datetime import datetime

import click as click
import log
import yaml


def append_path(current_path, item) -> str:
    return "".join([current_path.rstrip("/"), "/", str(item)  ])


def resolve_reference(reference: str, logger) -> dict:
    try:
        logger.debug(f'opening ref: {reference}')
        with open(reference, 'r') as stream:
            return yaml.load(stream, Loader=yaml.FullLoader)
    except IOError as e:
        print(f'error reading input file {reference} - error: {e}')
        return dict()


def resolve_references(data: dict, logger, path) -> dict:
    logger.debug(f"checking path : {path}")
    if not isinstance(data, dict):
        logger.debug(f"path: {path} the root path is not dict we give up")
        return data

    for key in data.keys():
        logger.debug(f'path: {path} key is: {key}')

        if key == '$ref' and not data[key].startswith('#/'):
            logger.debug(f'path: {path} found ref: {key} -> {data[key]}')
            ref_value = data[key]
            resolved_value = resolve_reference(data[key], logger)
            # Remove the reference
            logger.debug(f'removing ref: {key} -> {ref_value}')
            data.pop(key, None)
            # add resolved value
            logger.debug(f'replacing {ref_value} with {resolved_value}')
            data.update(resolved_value)
            # reprocess
            return resolve_references(data, logger, path)
        elif isinstance(data[key], dict):
            logger.debug(f'path: {path} key \'{key}\' value is a dict')
            newpath = append_path(path, key)
            resolved_val = resolve_references(data[key], logger, newpath)
            logger.debug(f"updating value at path : {path} for key : {key}")
            data[key] = resolved_val
        elif isinstance(data[key], list):
            logger.debug(f'path : {path} - key \'{key}\' value is a list')
            for i, s in enumerate(data[key]):
                newpath = append_path(path, str(i))
                logger.debug(f"path: ${newpath} replacing list item")
                data[key][i] = resolve_references(s, logger, newpath)

    return data


@click.command(context_settings={'help_option_names': ['-h', '--help']})
@click.option('--spec-input-file', is_flag=False, type=click.STRING,
              help='the input swagger specification file')
@click.option('--spec-output-file', is_flag=False, type=click.STRING,
              help='the input swagger specification file')
@click.option('--log-dir', is_flag=False, type=click.STRING, default='logs',
              help='the log output folder - default is the current folder')
def main(spec_input_file, spec_output_file, log_dir):
    if not os.path.isdir(log_dir):
        os.makedirs(log_dir)

    if spec_input_file is None:
        print(f'Input spec input file missing - set with the --spec-input-file option')
        return
    if spec_output_file is None:
        print(f'Input spec output file missing - set with the --spec-output-file option')
        return

    logger = log.setup_custom_logger('root', log_dir=log_dir, prefix='multi-file-swagger')
    logger.setLevel(logging.INFO)
    logger.info('log-dir     : {}'.format(log_dir))

    start_time = datetime.now()
    try:
        logger.info(f'Processing swagger spec file: {spec_input_file}')
        with open(spec_input_file, 'r') as stream:
            try:
                data = yaml.load(stream, Loader=yaml.FullLoader)
                data = resolve_references(data, logger, "/")
                with open(spec_output_file, 'w') as yaml_file:
                    yaml.dump(data, yaml_file, sort_keys=False)
            except yaml.YAMLError as exc:
                print(exc)
    except IOError as e:
        print(f'error reading input file {spec_input_file} - error: {e}')
        return
    logger.info(f'Wrote swagger spec file     : {spec_output_file}')
    logger.info(f'Finished - duration         : {datetime.now() - start_time}')


if __name__ == "__main__":
    sys.exit(main())
