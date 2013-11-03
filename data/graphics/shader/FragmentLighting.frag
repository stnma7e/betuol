#version 300

in vec4 diffuseColor;
in vec3 vertexNormal;
in vec3 modelSpacePosition;

out vec4 outputColor;

uniform vec3 modelSpaceLightPosition;
uniform vec4 lightIntensity;
uniform vec4 ambientIntensity;

void main()
{
    vec3 lightDir = normalize(modelSpacePosition - modelSpaceLightPosition);
    float cosAngIncidence = dot(normalize(vertexNormal), lightDir);
    cosAngIncidence = clamp(cosAngIncidence, 0, 1);
    
    /*vec4 fragColor = texture2D(sampler, uvPosition.st);*/
    vec4 fragColor = diffuseColor;
    outputColor = (fragColor * vec4(0.8,0.8,0.8,0.8) * cosAngIncidence) +
            (fragColor * vec4(0.2,0.2,0.2,0.2));
    /*outputColor = vec4(vertexNormal, 1) + (fragColor * ambientIntensity);*/
    /*outputColor = lightIntensity;*/
}
