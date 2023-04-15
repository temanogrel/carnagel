#!/usr/bin/env python

import argparse
import json
import os

import requests


class AphroditeInventory:
    def __init__(self):
        self.args = dict()

        self.parse_env_args()
        self.parse_cli_args()

        if self.args.list:
            data = self.generate_inventory()
        else:
            data = {}

        if self.args.pretty:
            print(json.dumps(data, sort_keys=True, indent=2))
        else:
            print(json.dumps(data))

    def parse_env_args(self):

        self.api_uri = os.environ.get('APHRODITE_API')
        if self.api_uri is None:
            raise ValueError('Missing environment variable APHRODITE_API')

        self.api_token = os.environ.get('APHRODITE_TOKEN')
        if self.api_token is None:
            raise ValueError('Missing environment variable APHRODITE_TOKEN')

    def parse_cli_args(self):
        parser = argparse.ArgumentParser(description='OpenStack Inventory Module')
        parser.add_argument('--private', action='store_true', help='Use private address for ansible host')
        parser.add_argument('--refresh', action='store_true', help='Refresh cached information')
        parser.add_argument('--debug', action='store_true', default=False, help='Enable debug output')
        parser.add_argument('--pretty', action='store_true', default=False, help='Pretty output')

        group = parser.add_mutually_exclusive_group(required=True)
        group.add_argument('--list', action='store_true', help='List active servers')
        group.add_argument('--host', help='List details about the specific host')

        self.args = parser.parse_args()

    def generate_inventory(self):

        # A few non-dynamic hosts
        data = {
            'rabbitmq': {
                'hosts': ['rabbitmq.sla.bz'],
                'vars': {
                    'ansible_user': os.environ.get('SERVER_SSH_USER'),
                    'ansible_ssh_pass': os.environ.get('SERVER_SSH_PASS')
                },
                'children': []
            },

            'aphrodite': {
                'hosts': ['aphrodite.sla.bz'],
                'vars': {
                    'ansible_user': os.environ.get('SERVER_SSH_USER'),
                    'ansible_ssh_pass': os.environ.get('SERVER_SSH_PASS')
                },
                'children': []
            },

            'scraper': {
                'hosts': ['rabbitmq-two.sla.bz'],
                'vars': {
                    'ansible_user': os.environ.get('SERVER_SSH_USER'),
                    'ansible_ssh_pass': os.environ.get('SERVER_SSH_PASS')
                },
                'children': []
            },

            'all': {
                'hosts': [],
                'vars': {
                    'ansible_user': os.environ.get('SERVER_SSH_USER'),
                    'ansible_ssh_pass': os.environ.get('SERVER_SSH_PASS')
                },
                'children': []
            }
        }

        for server in self._get_enabled_servers():
            if data.get(server['type']) is None:
                data[server['type']] = {
                    'hosts': [],
                    'vars': {},
                    'children': []
                }

            if data.get(server['network']) is None:
                data[server['network']] = {
                    'hosts': [],
                    'vars': {},
                    'children': []
                }

            data['all']['hosts'].append(server['hostname'])
            data[server['type']]['hosts'].append(server['hostname'])
            data[server['network']]['hosts'].append(server['hostname'])

        return data

    def _get_enabled_servers(self):
        """
        Ansible still does not support python3 WTF!?

        :return:
        """
        session = requests.Session()
        session.headers.update({'Authorization': 'server {}'.format(self.api_token)})

        return session.get(self.api_uri + '/servers', params=dict(enabled=1)).json()['data']


if __name__ == '__main__':
    AphroditeInventory()
