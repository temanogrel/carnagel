<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: struct.proto

namespace Pb;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * `Value` represents a dynamically typed value which can be either
 * null, a number, a string, a boolean, a recursive struct value, or a
 * list of values. A producer of value is expected to set one of that
 * variants, absence of any variant indicates an error.
 * The JSON representation for `Value` is JSON value.
 *
 * Generated from protobuf message <code>pb.Value</code>
 */
class Value extends \Google\Protobuf\Internal\Message
{
    protected $kind;

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type int $null_value
     *           Represents a null value.
     *     @type float $number_value
     *           Represents a double value.
     *     @type string $string_value
     *           Represents a string value.
     *     @type bool $bool_value
     *           Represents a boolean value.
     *     @type \Pb\Struct $struct_value
     *           Represents a structured value.
     *     @type \Pb\ListValue $list_value
     *           Represents a repeated `Value`.
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Struct::initOnce();
        parent::__construct($data);
    }

    /**
     * Represents a null value.
     *
     * Generated from protobuf field <code>.pb.NullValue null_value = 1;</code>
     * @return int
     */
    public function getNullValue()
    {
        return $this->readOneof(1);
    }

    /**
     * Represents a null value.
     *
     * Generated from protobuf field <code>.pb.NullValue null_value = 1;</code>
     * @param int $var
     * @return $this
     */
    public function setNullValue($var)
    {
        GPBUtil::checkEnum($var, \Pb\NullValue::class);
        $this->writeOneof(1, $var);

        return $this;
    }

    /**
     * Represents a double value.
     *
     * Generated from protobuf field <code>double number_value = 2;</code>
     * @return float
     */
    public function getNumberValue()
    {
        return $this->readOneof(2);
    }

    /**
     * Represents a double value.
     *
     * Generated from protobuf field <code>double number_value = 2;</code>
     * @param float $var
     * @return $this
     */
    public function setNumberValue($var)
    {
        GPBUtil::checkDouble($var);
        $this->writeOneof(2, $var);

        return $this;
    }

    /**
     * Represents a string value.
     *
     * Generated from protobuf field <code>string string_value = 3;</code>
     * @return string
     */
    public function getStringValue()
    {
        return $this->readOneof(3);
    }

    /**
     * Represents a string value.
     *
     * Generated from protobuf field <code>string string_value = 3;</code>
     * @param string $var
     * @return $this
     */
    public function setStringValue($var)
    {
        GPBUtil::checkString($var, True);
        $this->writeOneof(3, $var);

        return $this;
    }

    /**
     * Represents a boolean value.
     *
     * Generated from protobuf field <code>bool bool_value = 4;</code>
     * @return bool
     */
    public function getBoolValue()
    {
        return $this->readOneof(4);
    }

    /**
     * Represents a boolean value.
     *
     * Generated from protobuf field <code>bool bool_value = 4;</code>
     * @param bool $var
     * @return $this
     */
    public function setBoolValue($var)
    {
        GPBUtil::checkBool($var);
        $this->writeOneof(4, $var);

        return $this;
    }

    /**
     * Represents a structured value.
     *
     * Generated from protobuf field <code>.pb.Struct struct_value = 5;</code>
     * @return \Pb\Struct
     */
    public function getStructValue()
    {
        return $this->readOneof(5);
    }

    /**
     * Represents a structured value.
     *
     * Generated from protobuf field <code>.pb.Struct struct_value = 5;</code>
     * @param \Pb\Struct $var
     * @return $this
     */
    public function setStructValue($var)
    {
        GPBUtil::checkMessage($var, \Pb\Struct::class);
        $this->writeOneof(5, $var);

        return $this;
    }

    /**
     * Represents a repeated `Value`.
     *
     * Generated from protobuf field <code>.pb.ListValue list_value = 6;</code>
     * @return \Pb\ListValue
     */
    public function getListValue()
    {
        return $this->readOneof(6);
    }

    /**
     * Represents a repeated `Value`.
     *
     * Generated from protobuf field <code>.pb.ListValue list_value = 6;</code>
     * @param \Pb\ListValue $var
     * @return $this
     */
    public function setListValue($var)
    {
        GPBUtil::checkMessage($var, \Pb\ListValue::class);
        $this->writeOneof(6, $var);

        return $this;
    }

    /**
     * @return string
     */
    public function getKind()
    {
        return $this->whichOneof("kind");
    }

}

