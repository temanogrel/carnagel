<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: minion.proto

namespace Pb;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>pb.ServerIdentifier</code>
 */
class ServerIdentifier extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>string hostname = 1;</code>
     */
    private $hostname = '';

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type string $hostname
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Minion::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>string hostname = 1;</code>
     * @return string
     */
    public function getHostname()
    {
        return $this->hostname;
    }

    /**
     * Generated from protobuf field <code>string hostname = 1;</code>
     * @param string $var
     * @return $this
     */
    public function setHostname($var)
    {
        GPBUtil::checkString($var, True);
        $this->hostname = $var;

        return $this;
    }

}
