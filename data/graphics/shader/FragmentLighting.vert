#version 300

layout(location = 0) in vec3 position;
layout(location = 1) in vec3 normal;

out vec4 diffuseColor;
out vec3 vertexNormal;
out vec3 modelSpacePosition;

uniform mat4 modelToCamera;
uniform mat4 projection;

void main()
{
	gl_Position = vec4(position, 1.0) * modelToCamera * projection;

	modelSpacePosition = position;
	vertexNormal = normal;
	diffuseColor = vec4(1.0,1.0,1.0,1.0);
}
