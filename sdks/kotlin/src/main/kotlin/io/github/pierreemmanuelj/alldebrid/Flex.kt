package io.github.pierreemmanuelj.alldebrid

import kotlinx.serialization.KSerializer
import kotlinx.serialization.Serializable
import kotlinx.serialization.builtins.ListSerializer
import kotlinx.serialization.builtins.MapSerializer
import kotlinx.serialization.builtins.serializer
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import kotlinx.serialization.json.JsonArray
import kotlinx.serialization.json.JsonDecoder
import kotlinx.serialization.json.JsonElement
import kotlinx.serialization.json.JsonNull
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.JsonPrimitive
import kotlinx.serialization.json.JsonTransformingSerializer
import kotlinx.serialization.json.jsonPrimitive

/** A value the API sends either as a JSON string or as a JSON number, such as a delayed link ID. */
@Serializable(with = FlexStringSerializer::class)
@JvmInline
public value class FlexString(public val value: String) {
    /** The value as a number, when it is one. */
    public fun toLongOrNull(): Long? = value.toLongOrNull()

    override fun toString(): String = value
}

/** Decodes a JSON string or number into a [FlexString]. */
public object FlexStringSerializer : KSerializer<FlexString> {
    override val descriptor: SerialDescriptor = PrimitiveSerialDescriptor("FlexString", PrimitiveKind.STRING)

    override fun deserialize(decoder: Decoder): FlexString {
        val element = (decoder as JsonDecoder).decodeJsonElement()
        return FlexString(element.jsonPrimitive.content)
    }

    override fun serialize(encoder: Encoder, value: FlexString) {
        encoder.encodeString(value.value)
    }
}

/** Decodes one string or a list of strings into a list. */
public object StringsSerializer :
    JsonTransformingSerializer<List<String>>(ListSerializer(String.serializer())) {
    override fun transformDeserialize(element: JsonElement): JsonElement =
        if (element is JsonPrimitive) JsonArray(listOf(element)) else element
}

/** The magnets of a status: a list, or a single object when one magnet was asked for. */
public object MagnetsSerializer :
    JsonTransformingSerializer<List<Magnet>>(ListSerializer(Magnet.serializer())) {
    override fun transformDeserialize(element: JsonElement): JsonElement = when (element) {
        is JsonObject -> JsonArray(listOf(element))
        JsonNull -> JsonArray(emptyList())
        else -> element
    }
}

/** Quotas per host. The API sends an empty map as an empty list. */
public object QuotasSerializer :
    JsonTransformingSerializer<Map<String, Long>>(MapSerializer(String.serializer(), Long.serializer())) {
    override fun transformDeserialize(element: JsonElement): JsonElement {
        if (element is JsonArray) {
            require(element.isEmpty()) { "quotas sent as a non-empty list" }
            return JsonObject(emptyMap())
        }
        return element
    }
}
